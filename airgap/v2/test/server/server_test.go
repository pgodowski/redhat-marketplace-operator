// Copyright 2021 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/adminserver"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/fileretriever"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/filesender"
	v1 "github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/model"
	server "github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/internal/server"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/pkg/database"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/pkg/models"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var db database.Database
var dbName = "server.db"

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func runSetup() {
	//Initialize the mock connection and server
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	bs := &server.Server{}
	mockSenderServer := &server.FileSenderServer{Server: bs}
	mockRetrieverServer := &server.FileRetrieverServer{Server: bs}
	mockAdminServer := &server.AdminServer{Server: bs}

	//Initialize logger
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zapr, due to error: %v", err))
	}
	bs.Log = zapr.NewLogger(zapLog)

	//Create Sqlite Database
	gormDb, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error during creation of Database")
	}
	db.DB = gormDb
	db.Log = bs.Log

	//Create tables
	err = db.DB.AutoMigrate(&models.FileMetadata{}, &models.File{}, &models.Metadata{})
	if err != nil {
		log.Fatalf("Error during creation of Models: %v", err)
	}

	bs.FileStore = &db

	filesender.RegisterFileSenderServer(s, mockSenderServer)
	fileretriever.RegisterFileRetrieverServer(s, mockRetrieverServer)
	adminserver.RegisterAdminServerServer(s, mockAdminServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			if err.Error() != "closed" { //When lis of type (*bufconn.Listener) is closed, server doesn't have to panic.
				panic(err)
			}
		} else {
			log.Printf("Mock server started")
		}
	}()
}

func createClient() *grpc.ClientConn {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial bufnet: %v", err)
	}

	return conn
}

func shutdown(conn *grpc.ClientConn) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Error: Couldn't close Database: %v", err)
	}
	sqlDB.Close()
	conn.Close()
	os.Remove(dbName)
	lis.Close()
}

func TestFileSenderServer_UploadFile(t *testing.T) {
	//Initialize server
	runSetup()
	//Initialize client
	conn := createClient()
	client := filesender.NewFileSenderClient(conn)

	//Shutdown resources
	defer shutdown(conn)

	sampleData := make([]byte, 1024)
	tests := []struct {
		name    string
		info    *filesender.UploadFileRequest_Info
		data    *filesender.UploadFileRequest_ChunkData
		res     *filesender.UploadFileResponse
		errCode codes.Code
	}{
		{
			name: "invalid file upload request with nil file info/data",
			info: &filesender.UploadFileRequest_Info{
				Info: nil,
			},
			data: &filesender.UploadFileRequest_ChunkData{
				ChunkData: nil,
			},
			res:     &filesender.UploadFileResponse{},
			errCode: codes.Unknown,
		},
		{
			name: "valid file upload",
			info: &filesender.UploadFileRequest_Info{
				Info: &v1.FileInfo{
					FileId: &v1.FileID{
						Data: &v1.FileID_Name{
							Name: "test-file.zip",
						},
					},
					Size: 1024,
					Metadata: map[string]string{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
				},
			},
			data: &filesender.UploadFileRequest_ChunkData{
				ChunkData: sampleData,
			},
			res: &filesender.UploadFileResponse{
				Size: 1024,
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "test-file.zip",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uploadClient, err := client.UploadFile(context.Background())
			if err != nil {
				t.Errorf("error while invoking grpc method upload file for test:%v with err: %v", tt.name, err)
			}

			err = uploadClient.Send(&filesender.UploadFileRequest{
				Data: tt.info,
			})
			if err != nil {
				t.Errorf("error while sending metadata for test:%v with err: %v", tt.name, err)
			}

			err = uploadClient.Send(&filesender.UploadFileRequest{
				Data: tt.data,
			})
			if err != nil {
				t.Errorf("error while uploading byte stream for test:%v with err: %v", tt.name, err)
			}

			res, err := uploadClient.CloseAndRecv()

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Errorf("mismatched error codes: expected %v, received: %v for test: %v", tt.errCode, er.Code(), tt.name)
					}
				}
			}

			if res != nil {
				if res.Size != tt.res.Size {
					t.Errorf("sent:%v and recieved:%v size doesn't match for test: %v", tt.res.Size, res.Size, tt.name)
				}
				if res.FileId.GetName() != tt.info.Info.FileId.GetName() {
					t.Errorf("name of uploaded file and downloaded file name does not match: %v != %v for test: %v", res.FileId.GetName(), tt.info.Info.FileId.GetName(), tt.name)
				}
			}
		})
	}
}

func TestFileSenderServer_UpdateFileMetadata(t *testing.T) {
	//Initialize server
	runSetup()
	//Initialize client
	conn := createClient()

	//Shutdown resources
	defer shutdown(conn)

	populateDataset(conn, t)
	//Create a client for download
	client := filesender.NewFileSenderClient(conn)
	retClient := fileretriever.NewFileRetrieverClient(conn)

	//Shutdown resources
	defer shutdown(conn)

	tests := []struct {
		name    string
		req     *filesender.UpdateFileMetadataRequest
		resp    *filesender.UpdateFileMetadataResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "update metadata of file",
			req: &filesender.UpdateFileMetadataRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "reports.zip",
					},
				},
				Metadata: map[string]string{
					"version": "2",
					"type":    "report file",
				},
			},
			resp: &filesender.UpdateFileMetadataResponse{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "reports.zip",
					},
				},
			},
			errCode: codes.OK,
		},
		{
			name: "update metadata of file with same metadata",
			req: &filesender.UpdateFileMetadataRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "marketplace_report.zip",
					},
				},
				Metadata: map[string]string{
					"version": "2",
					"type":    "marketplace_report",
				},
			},
			resp:    nil,
			errCode: codes.Unknown,
			errMsg:  "connot update file metadata, as metadata of latest file and update request is same",
		},
		{
			name: "update metadata of non existing file",
			req: &filesender.UpdateFileMetadataRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "non exsiting",
					},
				},
				Metadata: map[string]string{
					"version": "2",
					"type":    "file type",
				},
			},
			resp:    nil,
			errCode: codes.Unknown,
			errMsg:  "connot update file metadata, as metadata of latest file and update request is same",
		},
	}

	for _, tt := range tests {
		_, err := client.UpdateFileMetadata(context.Background(), tt.req)
		if err != nil {
			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.errCode {
					t.Errorf("mismatched error codes: expected %v, received: %v for test: %v", tt.errCode, er.Code(), tt.name)
				}
			}
		}
		if tt.resp != nil {
			getResp, _ := retClient.GetFileMetadata(context.Background(), &fileretriever.GetFileMetadataRequest{FileId: tt.resp.GetFileId()})
			if !reflect.DeepEqual(tt.req.Metadata, getResp.GetInfo().GetMetadata()) {
				t.Errorf("metadata of file and metadata is update request does no match , Expected: %v got %v", tt.req.Metadata, getResp.GetInfo().GetMetadata())
			}
		}
	}
}

func TestFileRetrieverServer_DownloadFile(t *testing.T) {
	//Initialize server
	runSetup()
	//Initialize client
	conn := createClient()

	//Shutdown resources
	defer shutdown(conn)

	populateDataset(conn, t)
	//Create a client for download
	downloadClient := fileretriever.NewFileRetrieverClient(conn)

	tests := []struct {
		name    string
		dfr     *fileretriever.DownloadFileRequest
		size    uint32
		errCode codes.Code
	}{
		{
			name: "download an existing file on the server",
			dfr: &fileretriever.DownloadFileRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "reports.zip"},
				},
			},
			size:    2000,
			errCode: codes.OK,
		},
		{
			name: "download an existing file on the server and mark it for deletion",
			dfr: &fileretriever.DownloadFileRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "dosfstools"},
				},
				DeleteOnDownload: true,
			},
			size:    2000,
			errCode: codes.OK,
		},
		{
			name: "invalid download request for file that doesn't exist on the server",
			dfr: &fileretriever.DownloadFileRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "dontexist.zip"},
				},
			},
			size:    0,
			errCode: codes.InvalidArgument,
		},
		{
			name: "invalid download request for file with only whitespaces for the name",
			dfr: &fileretriever.DownloadFileRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "   "},
				},
			},
			size:    0,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := downloadClient.DownloadFile(context.Background(), tt.dfr)

			if err != nil {
				t.Errorf("error while invoking grpc method download file for test:%v with err: %v", tt.name, err)
			}

			var bs bytes.Buffer
			for {
				data, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					if er, ok := status.FromError(err); ok {
						if er.Code() != tt.errCode {
							t.Errorf("mismatched error codes: expected %v, received: %v for test: %v", tt.errCode, er.Code(), tt.name)
						}
					}
					break
				}
				bs.Write(data.GetChunkData())
			}

			if bs.Len() != int(tt.size) {
				t.Errorf("sent:%v and recieved:%v size doesn't match for test: %v", bs.Len(), int(tt.size), tt.name)
			}

			if tt.dfr.GetDeleteOnDownload() {
				lfr := &fileretriever.ListFileMetadataRequest{
					FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{
						{
							Key:      "provided_name",
							Operator: fileretriever.ListFileMetadataRequest_ListFileFilter_EQUAL,
							Value:    tt.dfr.FileId.GetName(),
						},
					},
					SortBy:              []*fileretriever.ListFileMetadataRequest_ListFileSort{},
					IncludeDeletedFiles: false,
				}
				listFileMetadataClient := fileretriever.NewFileRetrieverClient(conn)
				stream, err := listFileMetadataClient.ListFileMetadata(context.Background(), lfr)
				var data []*v1.FileInfo
				if err != nil {
					t.Errorf("error while invoking grpc method list file metadata for test:%v with err: %v", tt.name, err)
				}
				for {
					response, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						t.Errorf("error while fetching list of file")
						break
					}
					t.Logf("Received data: %v ", response.Results)
					data = append(data, response.GetResults())
				}
				if len(data) != 0 {
					t.Errorf("file marked for deletion should not be listed. expected: [] | got: %v ", data)
				}
			}
		})
	}
}

func TestFileRetrieverServer_ListFileMetadata(t *testing.T) {
	//Initialize server
	runSetup()
	//Initialize client
	conn := createClient()
	//Shutdown resources
	defer shutdown(conn)

	populateDataset(conn, t)
	listFileMetadataClient := fileretriever.NewFileRetrieverClient(conn)

	tests := []struct {
		name    string
		lfr     *fileretriever.ListFileMetadataRequest
		res_len int
		errCode codes.Code
	}{
		{
			name: "fetch list of all file by passing empty filter array",
			lfr: &fileretriever.ListFileMetadataRequest{
				FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{},
				SortBy:   []*fileretriever.ListFileMetadataRequest_ListFileSort{},
			},
			res_len: 4,
			errCode: codes.OK,
		},
		{
			name: "fetch file list based on filter operation",
			lfr: &fileretriever.ListFileMetadataRequest{
				FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{
					{
						Key:      "description",
						Operator: fileretriever.ListFileMetadataRequest_ListFileFilter_CONTAINS,
						Value:    "filesystem utilities",
					},
				},
				SortBy: []*fileretriever.ListFileMetadataRequest_ListFileSort{},
			},
			res_len: 1,
			errCode: codes.OK,
		},
		{
			name: "fetch file marked for deletion",
			lfr: &fileretriever.ListFileMetadataRequest{
				FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{
					{
						Key:      "provided_name",
						Operator: fileretriever.ListFileMetadataRequest_ListFileFilter_CONTAINS,
						Value:    "delete",
					},
				},
				SortBy:              []*fileretriever.ListFileMetadataRequest_ListFileSort{},
				IncludeDeletedFiles: true,
			},
			res_len: 1,
			errCode: codes.OK,
		},
		{
			name: "empty values in filter operation",
			lfr: &fileretriever.ListFileMetadataRequest{
				FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{
					{},
				},
				SortBy: []*fileretriever.ListFileMetadataRequest_ListFileSort{},
			},
			res_len: 0,
			errCode: codes.InvalidArgument,
		},
		{
			name: "empty values in sort operation",
			lfr: &fileretriever.ListFileMetadataRequest{
				FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{},
				SortBy: []*fileretriever.ListFileMetadataRequest_ListFileSort{
					{},
				},
			},
			res_len: 0,
			errCode: codes.InvalidArgument,
		},
		{
			name: "empty key/value for filter operation",
			lfr: &fileretriever.ListFileMetadataRequest{
				FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{
					{
						Key:      "     ",
						Operator: fileretriever.ListFileMetadataRequest_ListFileFilter_CONTAINS,
						Value:    "",
					},
				},
				SortBy: []*fileretriever.ListFileMetadataRequest_ListFileSort{},
			},
			res_len: 0,
			errCode: codes.InvalidArgument,
		},
		{
			name: "empty sort key for sort operation",
			lfr: &fileretriever.ListFileMetadataRequest{
				FilterBy: []*fileretriever.ListFileMetadataRequest_ListFileFilter{},
				SortBy: []*fileretriever.ListFileMetadataRequest_ListFileSort{
					{
						Key:       "  ",
						SortOrder: fileretriever.ListFileMetadataRequest_ListFileSort_DESC,
					},
				},
			},
			res_len: 0,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := listFileMetadataClient.ListFileMetadata(context.Background(), tt.lfr)
			var data []*v1.FileInfo
			if err != nil {
				t.Errorf("error while invoking grpc method list file metadata for test:%v with err: %v", tt.name, err)
			}
			for {
				response, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					if er, ok := status.FromError(err); ok {
						if er.Code() != tt.errCode {
							t.Errorf("mismatched error codes: expected %v, received: %v, details: %v | for test: %v",
								tt.errCode, er.Code(), er.Message(), tt.name)
						}
					}
					break
				}
				t.Logf("Received data: %v ", response.Results)
				data = append(data, response.GetResults())
			}
			if len(data) != tt.res_len {
				t.Errorf("requested data and received data doesn't match for test: %v ", tt.name)
			}
		})
	}
}

func TestFileRetrieverServer_GetFileMetadata(t *testing.T) {
	//Initialize server
	runSetup()
	//Initialize client
	conn := createClient()

	//Shutdown resources
	defer shutdown(conn)

	populateDataset(conn, t)

	//Create a client for download
	getFileMetadaClient := fileretriever.NewFileRetrieverClient(conn)

	tests := []struct {
		name    string
		dfr     *fileretriever.GetFileMetadataRequest
		size    uint32
		errCode codes.Code
	}{
		{
			name: "get file metadata of an existing file on the server",
			dfr: &fileretriever.GetFileMetadataRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "reports.zip"},
				},
			},
			size:    2000,
			errCode: codes.OK,
		},
		{
			name: "invalid get file metadata request for file that doesn't exist on the server",
			dfr: &fileretriever.GetFileMetadataRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "dontexist.zip"},
				},
			},
			size:    0,
			errCode: codes.InvalidArgument,
		},
		{
			name: "invalid get file metadata request for file with only whitespaces for the name",
			dfr: &fileretriever.GetFileMetadataRequest{
				FileId: &v1.FileID{
					Data: &v1.FileID_Name{
						Name: "   "},
				},
			},
			size:    0,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := getFileMetadaClient.GetFileMetadata(context.Background(), tt.dfr)

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Errorf("mismatched error codes: expected %v, received: %v for test: %v", tt.errCode, er.Code(), tt.name)
					}
				}
			}

			md := stream.GetInfo()

			if md != nil {
				if int(md.GetSize()) != int(tt.size) {
					t.Errorf("sent:%v and recieved:%v size doesn't match for test: %v", int(tt.size), md.GetSize(), tt.name)
				}
			}
		})
	}
}

func TestAdminServer_CleanTombstones(t *testing.T) {
	//Initialize server
	runSetup()
	//Initialize client
	conn := createClient()

	//Shutdown resources
	defer shutdown(conn)

	populateDataset(conn, t)

	//Create a client for download
	client := adminserver.NewAdminServerClient(conn)

	tests := []struct {
		name    string
		req     *adminserver.CleanTombstonesRequest
		resp    *adminserver.CleanTombstonesResponse
		errCode codes.Code
	}{
		{
			name: "clean file content",
			req: &adminserver.CleanTombstonesRequest{
				Before:   timestamppb.Now(),
				PurgeAll: false,
			},
			resp: &adminserver.CleanTombstonesResponse{
				TombstonesCleaned: 1,
				Files: []*v1.FileID{
					{
						Data: &v1.FileID_Name{Name: "delete.txt"},
					},
				},
			},
		},
		{
			name: "delete file record",
			req: &adminserver.CleanTombstonesRequest{
				Before:   &timestamppb.Timestamp{Seconds: time.Now().Unix()},
				PurgeAll: true,
			},
			resp: &adminserver.CleanTombstonesResponse{
				TombstonesCleaned: 1,
				Files: []*v1.FileID{
					{
						Data: &v1.FileID_Name{Name: "delete.txt"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.CleanTombstones(context.Background(), tt.req)
			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Errorf("mismatched error codes: expected %v, received: %v for test: %v", tt.errCode, er.Code(), tt.name)
						t.Log(err)
					}
				}
			}
			if resp.GetTombstonesCleaned() != tt.resp.GetTombstonesCleaned() {
				t.Errorf("mismatched files cleaned: expected %v, received: %v", tt.resp.Files, resp.GetFiles())
			} else if !reflect.DeepEqual(resp.GetFiles(), tt.resp.GetFiles()) {
				t.Errorf("mismatched files cleaned: expected %v, received: %v", tt.resp.Files, resp.GetFiles())
			}
		})
	}

}

// populateDataset populates database with the files needed for testing
func populateDataset(conn *grpc.ClientConn, t *testing.T) {
	deleteFID := &v1.FileID{
		Data: &v1.FileID_Name{
			Name: "delete.txt",
		}}
	files := []v1.FileInfo{
		{
			FileId: &v1.FileID{
				Data: &v1.FileID_Name{
					Name: "reports.zip",
				},
			},
			Size:            2000,
			Compression:     true,
			CompressionType: "gzip",
			Metadata: map[string]string{
				"version": "2",
				"type":    "report",
			},
		},
		{
			FileId: &v1.FileID{
				Data: &v1.FileID_Name{
					Name: "marketplace_report.zip",
				},
			},
			Size:            200,
			Compression:     true,
			CompressionType: "gzip",
			Metadata: map[string]string{
				"version": "2",
				"type":    "marketplace_report",
			},
		},
		{
			FileId: &v1.FileID{
				Data: &v1.FileID_Name{
					Name: "dosfstools",
				},
			},
			Size:            2000,
			Compression:     true,
			CompressionType: "gzip",
			Metadata: map[string]string{
				"version":     "latest",
				"description": "DOS filesystem utilities",
			},
		},
		{
			FileId: &v1.FileID{
				Data: &v1.FileID_Name{
					Name: "dosbox",
				},
			},
			Size:            1500,
			Compression:     true,
			CompressionType: "gzip",
			Metadata: map[string]string{
				"version":     "4.3",
				"description": "Emulator with builtin DOS for running DOS Games",
			},
		},
		{
			FileId:      deleteFID,
			Size:        1500,
			Compression: false,
			Metadata: map[string]string{
				"description": "file marked for deletion",
			},
		},
	}
	uploadClient := filesender.NewFileSenderClient(conn)

	// Upload files to mock server
	for i := range files {
		clientStream, err := uploadClient.UploadFile(context.Background())
		if err != nil {
			t.Fatalf("Error: During call of client.UploadFile: %v", err)
		}

		//Upload metadata
		err = clientStream.Send(&filesender.UploadFileRequest{
			Data: &filesender.UploadFileRequest_Info{
				Info: &files[i],
			},
		})

		if err != nil {
			t.Fatalf("Error: during sending metadata: %v", err)
		}

		//Upload chunk data
		bs := make([]byte, files[i].GetSize())
		request := filesender.UploadFileRequest{
			Data: &filesender.UploadFileRequest_ChunkData{
				ChunkData: bs,
			},
		}
		clientStream.Send(&request)

		res, err := clientStream.CloseAndRecv()
		if err != nil {
			t.Fatalf("Error: during stream close and recieve: %v", err)
		}
		t.Logf("Received response: %v", res)
	}

	// Mark File for deletion
	req := &fileretriever.DownloadFileRequest{
		FileId:           deleteFID,
		DeleteOnDownload: true,
	}
	dc := fileretriever.NewFileRetrieverClient(conn)
	_, err := dc.DownloadFile(context.Background(), req)
	if err != nil {
		t.Fatalf("Error: during delete on download request : %v", err)
	}
	time.Sleep(1 * time.Second)
}
