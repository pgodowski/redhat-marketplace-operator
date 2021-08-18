package ci

import (
	"tool/file"
	"encoding/yaml"
)

command: genworkflows: task: {
	for w in workflows {
		"\(w.file)": file.Create & {
			filename: w.file
			contents: """
				# Generated by internal/ci/ci_tool.cue; do not edit
				\(yaml.Marshal(w.schema))
				"""
		}
	}
}

command: genscripts: task: {
	for w in scripts {
		"\(w.file)": file.Create & {
			filename:    w.file
			permissions: 0o755
			contents:    """
			#!/bin/bash
			# Generated by internal/ci/ci_tool.cue; do not edit
			\(w.script.result)
			"""
		}
	}
}

command: gentravis: task: {
	for w in travis {
		"\(w.file)": file.Create & {
			filename: w.file
			contents: """
				# Generated by internal/ci/ci_tool.cue; do not edit

				\(yaml.Marshal(w.schema))
				"""
		}
	}
}
