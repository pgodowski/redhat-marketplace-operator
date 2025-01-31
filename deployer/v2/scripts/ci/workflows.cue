package ci

import (
  "github.com/SchemaStore/schemastore/src/schemas/json"
  encjson "encoding/json"
)

workflowsDir: *"./" | string @tag(workflowsDir)

workflows: [...{file: string, schema: (json.#Workflow & {})}]
workflows: [
  {
    file: "deploy_cue.yml"
    schema: deploy
  },
]

deploy: _#bashWorkflow & {
  name: "Test & Deploy Image"
  on: {
		push: {
			branches: [
        "master",
        "release/**",
        "hotfix/**",
        "develop",
        "feature/**",
        "bugfix/**",
      ]
		}
	}

  env: {
    pushImage: true
  }

  jobs: {
    "test-unit": {
      name: ""
      "runs-on": _#linuxMachine
      needs: ["build"]
      env: {
        "IMAGE_REGISTRY": "${{ matrix.registry }}"
        "OPERATOR_IMAGE": "${{ matrix.registry }}/redhat-marketplace-operator:${{ needs.build.outputs.dockertag }}"
        "OPERATOR_IMAGE_TAG": "${{ needs.build.outputs.dockertag }}"
        "TAG": "${{ matrix.registry }}/redhat-marketplace-operator:${{ needs.build.outputs.dockertag}}"
      }
      steps: [
				_#checkoutCode,
        _#setBranchPrefixForDev,
      ]
    }
  }
}

test: _#bashWorkflow & {
	name: "Test"
	on: {
		push: {
			branches: ["**"] // any branch (including '/' namespaced branches)
			"tags-ignore": ["v*"]
		}
	}

	jobs: {
		start: {
			"runs-on": _#linuxMachine
			if:        "${{ \(_#isCLCITestBranch) }}"
			steps: [
				_#writeCookiesFile,
				_#startCLBuild,
			]
		}
		test: {
			strategy:  _#testStrategy
			"runs-on": "${{ matrix.os }}"
			steps: [
				_#writeCookiesFile,
				_#installGo,
				_#checkoutCode,
				_#cacheGoModules,
				_#goGenerate,
				_#goTest,
				_#goTestRace,
				_#goReleaseCheck,
				_#checkGitClean,
				_#pullThroughProxy,
				_#failCLBuild,
			]
		}
		mark_ci_success: {
			"runs-on": _#linuxMachine
			if:        "${{ \(_#isCLCITestBranch) }}"
			needs:     "test"
			steps: [
				_#writeCookiesFile,
				_#passCLBuild,
			]
		}
		delete_build_branch: {
			"runs-on": _#linuxMachine
			if:        "${{ \(_#isCLCITestBranch) && always() }}"
			needs:     "test"
			steps: [
				_#step & {
					run: """
						\(_#tempCueckooGitDir)
						git push https://github.com/cuelang/cue :${GITHUB_REF#\(_#branchRefPrefix)}
						"""
				},
			]
		}
	}

	// _#isCLCITestBranch is an expression that evaluates to true
	// if the job is running as a result of a CL triggered CI build
	_#isCLCITestBranch: "startsWith(github.ref, '\(_#branchRefPrefix)ci/')"

	// _#isMaster is an expression that evaluates to true if the
	// job is running as a result of a master commit push
	_#isMaster: "github.ref == '\(_#branchRefPrefix)master'"

	_#pullThroughProxy: _#step & {
		name: "Pull this commit through the proxy on master"
		run: """
			v=$(git rev-parse HEAD)
			cd $(mktemp -d)
			go mod init mod.com
			GOPROXY=https://proxy.golang.org go get -d cuelang.org/go@$v
			"""
		if: "${{ \(_#isMaster) }}"
	}

	_#startCLBuild: _#step & {
		name: "Update Gerrit CL message with starting message"
		run:  (_#gerrit._#setCodeReview & {
			#args: message: "Started the build... see progress at ${{ github.event.repository.html_url }}/actions/runs/${{ github.run_id }}"
		}).res
	}

	_#failCLBuild: _#step & {
		if:   "${{ \(_#isCLCITestBranch) && failure() }}"
		name: "Post any failures for this matrix entry"
		run:  (_#gerrit._#setCodeReview & {
			#args: {
				message: "Build failed for ${{ runner.os }}-${{ matrix.go-version }}; see ${{ github.event.repository.html_url }}/actions/runs/${{ github.run_id }} for more details"
				labels: {
					"Code-Review": -1
				}
			}
		}).res
	}

	_#passCLBuild: _#step & {
		name: "Update Gerrit CL message with success message"
		run:  (_#gerrit._#setCodeReview & {
			#args: {
				message: "Build succeeded for ${{ github.event.repository.html_url }}/actions/runs/${{ github.run_id }}"
				labels: {
					"Code-Review": 1
				}
			}
		}).res
	}

	_#gerrit: {
		// _#setCodeReview assumes that it is invoked from a job where
		// _#isCLCITestBranch is true
		_#setCodeReview: {
			#args: {
				message: string
				labels?: {
					"Code-Review": int
				}
			}
			res: #"""
			curl -f -s -H "Content-Type: application/json" --request POST --data '\#(encjson.Marshal(#args))' -b ~/.gitcookies https://cue-review.googlesource.com/a/changes/$(basename $(dirname $GITHUB_REF))/revisions/$(basename $GITHUB_REF)/review
			"""#
		}
	}
}

test_dispatch: _#bashWorkflow & {

	name: "Test Dispatch"
	on: ["repository_dispatch"]
	jobs: {
		start: {
			if:        "${{ startsWith(github.event.action, 'Build for refs/changes/') }}"
			"runs-on": _#linuxMachine
			steps: [
				_#step & {
					name: "Checkout ref"
					run:  """
						\(_#tempCueckooGitDir)
						git fetch https://cue-review.googlesource.com/cue ${{ github.event.client_payload.ref }}
						git checkout -b ci/${{ github.event.client_payload.changeID }}/${{ github.event.client_payload.commit }} FETCH_HEAD
						git push https://github.com/cuelang/cue ci/${{ github.event.client_payload.changeID }}/${{ github.event.client_payload.commit }}
						"""
				},
			]
		}
	}
}

release: _#bashWorkflow & {

	name: "Release"
	on: push: tags: ["v*"]
	jobs: {
		goreleaser: {
			"runs-on": _#linuxMachine
			steps: [{
				name: "Checkout code"
				uses: "actions/checkout@v3"
			}, {
				name: "Unshallow" // required for the changelog to work correctly.
				run:  "git fetch --prune --unshallow"
			}, {
				name: "Run GoReleaser"
				env: GITHUB_TOKEN: "${{ secrets.ACTIONS_GITHUB_TOKEN }}"
				uses: "docker://goreleaser/goreleaser:latest"
				with: args: "release --rm-dist"
			}]
		}
		docker: {
			name:      "docker"
			"runs-on": _#linuxMachine
			steps: [{
				name: "Check out the repo"
				uses: "actions/checkout@v3"
			}, {
				name: "Set version environment"
				run: """
					CUE_VERSION=$(echo ${GITHUB_REF##refs/tags/v})
					echo \"CUE_VERSION=$CUE_VERSION\"
					echo \"CUE_VERSION=$(echo $CUE_VERSION)\" >> $GITHUB_ENV
					"""
			}, {
				name: "Push to Docker Hub"
				env: {
					DOCKER_BUILDKIT: 1
					GOLANG_VERSION:  1.14
					CUE_VERSION:     "${{ env.CUE_VERSION }}"
				}
				uses: "docker/build-push-action@v2"
				with: {
					tags:           "${{ env.CUE_VERSION }},latest"
					repository:     "cuelang/cue"
					username:       "${{ secrets.DOCKER_USERNAME }}"
					password:       "${{ secrets.DOCKER_PASSWORD }}"
					tag_with_ref:   false
					tag_with_sha:   false
					target:         "cue"
					always_pull:    true
					build_args:     "GOLANG_VERSION=${{ env.GOLANG_VERSION }},CUE_VERSION=v${{ env.CUE_VERSION }}"
					add_git_labels: true
				}
			}]
		}
	}
}

rebuild_tip_cuelang_org: _#bashWorkflow & {

	name: "Push to tip"
	on: push: branches: ["master"]
	jobs: push: {
		"runs-on": _#linuxMachine
		steps: [{
			name: "Rebuild tip.cuelang.org"
			run:  "curl -f -X POST -d {} https://api.netlify.com/build_hooks/${{ secrets.CuelangOrgTipRebuildHook }}"
		}]
	}
}

_#bashWorkflow: json.#Workflow & {
	jobs: [string]: defaults: run: shell: "bash"
}

// TODO: drop when cuelang.org/issue/390 is fixed.
// Declare definitions for sub-schemas
_#job:  ((json.#Workflow & {}).jobs & {x: _}).x
_#step: ((_#job & {steps:                 _}).steps & [_])[0]

// We need at least go1.14 for code generation
_#codeGenGo: "1.14.9"

_#linuxMachine:   "ubuntu-20.04"
_#macosMachine:   "macos-10.15"
_#windowsMachine: "windows-2019"

_#testStrategy: {
	"fail-fast": false
	matrix: {
		// Use a stable version of 1.14.x for go generate
		"go-version": ["1.13.x", _#codeGenGo, "1.15.x"]
		os: [_#linuxMachine, _#macosMachine, _#windowsMachine]
	}
}

_#cancelPreviousRun: _#step & {
  name: "Cancel Previous Run"
  uses: "styfle/cancel-workflow-action@0.11.0"
  with: "access_token": "${{ github.token }}"
}

_#installGo: _#step & {
	name: "Install Go"
	uses: "actions/setup-go@v4"
	with: {
		"go-version": "${{ matrix.go-version }}"
		cache: false
	}
}

_#checkoutCode: _#step & {
	name: "Checkout code"
	uses: "actions/checkout@v3"
}

_#cacheGoModules: _#step & {
	name: "Cache Go modules"
	uses: "actions/cache@v3"
	with: {
		path: "~/go/pkg/mod"
		key:  "${{ runner.os }}-${{ matrix.go-version }}-go-${{ hashFiles('**/go.sum') }}"
		"restore-keys": """
			${{ runner.os }}-${{ matrix.go-version }}-go-
			"""
	}
}

_#goGenerate: _#step & {
	name: "Generate"
	run:  "go generate ./..."
	// The Go version corresponds to the precise version specified in
	// the matrix. Skip windows for now until we work out why re-gen is flaky
	if: "matrix.go-version == '\(_#codeGenGo)' && matrix.os != '\(_#windowsMachine)'"
}

_#goTest: _#step & {
	name: "Test"
	run:  "go test ./..."
}

_#goTestRace: _#step & {
	name: "Test with -race"
	run:  "go test -race ./..."
}

_#goReleaseCheck: _#step & {
	name: "gorelease check"
	run:  "go run golang.org/x/exp/cmd/gorelease"
}

_#loadGitTagPushed: _#step & {
  name: "Get if gittag is pushed"
  id: "tag"
  run: """
  VERSION=$(make current-version)
  RESULT=$(git tag --list | grep -E "$VERSION")
  IS_TAGGED=false
  if [ "$RESULT" != "" ] ; then
    IS_TAGGED=true
  """
}

_#branchRefPrefix: "refs/heads/"

_#tempCueckooGitDir: """
	mkdir tmpgit
	cd tmpgit
	git init
	git config user.name cueckoo
	git config user.email cueckoo@gmail.com
	git config http.https://github.com/.extraheader "AUTHORIZATION: basic $(echo -n cueckoo:${{ secrets.CUECKOO_GITHUB_PAT }} | base64)"
	"""

_#setBranchPrefixForDev: (_#vars._#setBranchPrefix & {
  #args: {
    eventName: "push"
    branch: "refs/heads/develop"
    tagPrefix: "dev-"
  }
}).res

	_#vars: {
		// _#setBranchPrefix will set the branch prefix vars
		_#setBranchPrefix: {
			#args: {
				eventName: string
        branch: string
        tagPrefix: string
				quayExpiration?: string
			}
      res: _#step & {
        if: "github.event_name == '#args.eventName' && github.ref == '#(args.branch)'"
        run: """
        echo "TAGPREFIX=#(args.tagPrefix)" >> $GITHUB_ENV
        if [ "#(args.quayExpiration)" != "" ]; then
          echo "QUAY_EXPIRATION=$(args.quayExpiration)" >> $GITHUB_ENV
        fi
        """
      }
		}
	}
}


//       - name: Set branch prefix for dev
//         if: github.event_name == 'push' && github.ref == 'refs/heads/develop'
//         run: |


//       - name: Set branch prefix for fix
//         if: github.event_name == 'push' && startsWith(github.ref,'refs/heads/bugfix/')
//         run: |
//           NAME=$(echo "${{ github.ref }}" | sed 's/refs\/heads\/bugfix\///')
//           echo "TAGPREFIX=bugfix-${NAME}-" >> $GITHUB_ENV
//           echo "QUAY_EXPIRATION=1w" >> $GITHUB_ENV

//       - name: Set branch prefix for feat
//         if: github.event_name == 'push' && startsWith(github.ref,'refs/heads/feature/')
//         run: |
//           NAME=$(echo "${{ github.ref }}" | sed 's/refs\/heads\/feature\///')
//           echo "TAGPREFIX=feat-${NAME}-" >> $GITHUB_ENV
// echo "QUAY_EXPIRATION=1w" >> $GITHUB_ENV
