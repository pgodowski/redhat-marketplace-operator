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

package server

import (
	"context"
	"fmt"
	"io"

	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/filesender"
	v1 "github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileSenderServer struct {
	*Server
	filesender.UnimplementedFileSenderServer
}

type fileSenderuploadFileStream interface {
	SendAndClose(*filesender.UploadFileResponse) error
	Recv() (*filesender.UploadFileRequest, error)
}

func (frs *FileSenderServer) UploadFile(stream filesender.FileSender_UploadFileServer) error {
	return frs.uploadFile(stream)
}

// UploadFile allows a file to be uploaded and saved in the database
func (frs *FileSenderServer) uploadFile(stream fileSenderuploadFileStream) error {
	var bs []byte
	var finfo *v1.FileInfo
	var fid *v1.FileID
	var size uint32

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				frs.Log.Info("Stream end", "total bytes received", len(bs))
				// Attempt to save file in database
				err := frs.FileStore.SaveFile(finfo, bs)
				if err != nil {
					return status.Errorf(
						codes.Unknown,
						fmt.Sprintf("Failed to save file in database: %v", err),
					)
				}

				// Prepare response on save and close stream
				res := &filesender.UploadFileResponse{
					FileId: fid,
					Size:   size,
				}

				return stream.SendAndClose(res)
			}

			frs.Log.Error(err, "Oops, something went wrong!")
			return status.Errorf(
				codes.Unknown,
				fmt.Sprintf("Error while processing stream, details: %v", err),
			)
		}

		b := req.GetChunkData()
		if b != nil {
			if bs == nil {
				bs = b
			} else {
				bs = append(bs, b...)
			}
		}

		if req.GetInfo() != nil {
			finfo = req.GetInfo()
			fid = finfo.GetFileId()
			size = finfo.GetSize()
		}
	}
}

// UpdateFileMetadata allows to update metadata of file saved in the databse
func (frs *FileSenderServer) UpdateFileMetadata(ctx context.Context, in *filesender.UpdateFileMetadataRequest) (*filesender.UpdateFileMetadataResponse, error) {
	err := frs.FileStore.UpdateFileMetadata(in.GetFileId(), in.GetMetadata())
	if err != nil {
		return nil, status.Errorf(
			codes.Unknown,
			fmt.Sprintf("Failed to update: %v", err),
		)
	}

	response := filesender.UpdateFileMetadataResponse{
		FileId: in.GetFileId(),
	}

	return &response, nil
}
