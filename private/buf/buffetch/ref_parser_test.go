// Copyright 2020-2022 Buf Technologies, Inc.
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

package buffetch

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bufbuild/buf/private/buf/buffetch/internal"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduletesting"
	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/git"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetParsedRefSuccess(t *testing.T) {
	// This allows us to test an os-agnostic root directory
	root, err := filepath.Abs("/")
	require.NoError(t, err)

	// This lets us test an os-agnostic absolute path
	absPath, err := filepath.Abs("/foo/bar/..")
	require.NoError(t, err)
	expectedAbsDir, err := filepath.Abs("/foo")
	require.NoError(t, err)

	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedDirRef(
			formatDir,
			"path/to/some/dir",
		),
		"path/to/some/dir",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedDirRef(
			formatDir,
			".",
		),
		".",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedDirRef(
			formatDir,
			normalpath.Normalize(root),
		),
		root,
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedDirRef(
			formatDir,
			".",
		),
		"foo/..",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedDirRef(
			formatDir,
			"../foo",
		),
		"../foo",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedDirRef(
			formatDir,
			normalpath.Normalize(expectedAbsDir),
		),
		absPath,
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeNone,
			0,
			"",
		),
		"path/to/file.tar",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"/path/to/file.tar",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeNone,
			0,
			"",
		),
		"file:///path/to/file.tar",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeNone,
			1,
			"",
		),
		"path/to/file.tar#strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar.gz",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeGzip,
			0,
			"",
		),
		"path/to/file.tar.gz",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar.gz",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeGzip,
			1,
			"",
		),
		"path/to/file.tar.gz#strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tgz",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeGzip,
			0,
			"",
		),
		"path/to/file.tgz",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tgz",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeGzip,
			1,
			"",
		),
		"path/to/file.tgz#strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar",
			internal.FileSchemeHTTP,
			internal.ArchiveTypeTar,
			internal.CompressionTypeNone,
			0,
			"",
		),
		"http://path/to/file.tar",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar",
			internal.FileSchemeHTTPS,
			internal.ArchiveTypeTar,
			internal.CompressionTypeNone,
			0,
			"",
		),
		"https://path/to/file.tar",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatZip,
			"path/to/file.zip",
			internal.FileSchemeLocal,
			internal.ArchiveTypeZip,
			internal.CompressionTypeNone,
			0,
			"",
		),
		"path/to/file.zip",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatZip,
			"/path/to/file.zip",
			internal.FileSchemeLocal,
			internal.ArchiveTypeZip,
			internal.CompressionTypeNone,
			0,
			"",
		),
		"file:///path/to/file.zip",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatZip,
			"path/to/file.zip",
			internal.FileSchemeLocal,
			internal.ArchiveTypeZip,
			internal.CompressionTypeNone,
			1,
			"",
		),
		"path/to/file.zip#strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeLocal,
			nil,
			false,
			1,
			"",
		),
		"path/to/dir.git",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeLocal,
			nil,
			false,
			40,
			"",
		),
		"path/to/dir.git#depth=40",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeLocal,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"path/to/dir.git#branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"/path/to/dir.git",
			internal.GitSchemeLocal,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"file:///path/to/dir.git#branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeLocal,
			git.NewTagName("v1.0.0"),
			false,
			1,
			"",
		),
		"path/to/dir.git#tag=v1.0.0",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"hello.com/path/to/dir.git",
			internal.GitSchemeHTTP,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"http://hello.com/path/to/dir.git#branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"hello.com/path/to/dir.git",
			internal.GitSchemeHTTPS,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"https://hello.com/path/to/dir.git#branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"user@hello.com:path/to/dir.git",
			internal.GitSchemeSSH,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"ssh://user@hello.com:path/to/dir.git#branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"user@hello.com:path/to/dir.git",
			internal.GitSchemeSSH,
			git.NewRefName("refs/remotes/origin/HEAD"),
			false,
			50,
			"",
		),
		"ssh://user@hello.com:path/to/dir.git#ref=refs/remotes/origin/HEAD",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"user@hello.com:path/to/dir.git",
			internal.GitSchemeSSH,
			git.NewRefNameWithBranch("refs/remotes/origin/HEAD", "main"),
			false,
			50,
			"",
		),
		"ssh://user@hello.com:path/to/dir.git#ref=refs/remotes/origin/HEAD,branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"user@hello.com:path/to/dir.git",
			internal.GitSchemeSSH,
			git.NewRefName("refs/remotes/origin/HEAD"),
			false,
			10,
			"",
		),
		"ssh://user@hello.com:path/to/dir.git#ref=refs/remotes/origin/HEAD,depth=10",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"user@hello.com:path/to/dir.git",
			internal.GitSchemeSSH,
			git.NewRefNameWithBranch("refs/remotes/origin/HEAD", "main"),
			false,
			10,
			"",
		),
		"ssh://user@hello.com:path/to/dir.git#ref=refs/remotes/origin/HEAD,branch=main,depth=10",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeLocal,
			nil,
			false,
			1,
			"foo/bar",
		),
		"path/to/dir.git#subdir=foo/bar",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeLocal,
			nil,
			false,
			1,
			"",
		),
		"path/to/dir.git#subdir=.",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeLocal,
			nil,
			false,
			1,
			"",
		),
		"path/to/dir.git#subdir=foo/..",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"user@hello.com:path/to/dir.git",
			internal.GitSchemeGit,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"git://user@hello.com:path/to/dir.git#branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir.git",
			internal.GitSchemeGit,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"git://path/to/dir.git#branch=main",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"path/to/file.bin",
			internal.FileSchemeLocal,
			internal.CompressionTypeNone,
		),
		"path/to/file.bin",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"path/to/file.bin.gz",
			internal.FileSchemeLocal,
			internal.CompressionTypeGzip,
		),
		"path/to/file.bin.gz",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatJSON,
			"path/to/file.json",
			internal.FileSchemeLocal,
			internal.CompressionTypeNone,
		),
		"path/to/file.json",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatJSON,
			"path/to/file.json.gz",
			internal.FileSchemeLocal,
			internal.CompressionTypeGzip,
		),
		"path/to/file.json.gz",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatJSON,
			"path/to/file.json.gz",
			internal.FileSchemeLocal,
			internal.CompressionTypeNone,
		),
		"path/to/file.json.gz#compression=none",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatJSON,
			"path/to/file.json.gz",
			internal.FileSchemeLocal,
			internal.CompressionTypeGzip,
		),
		"path/to/file.json.gz#compression=gzip",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"",
			internal.FileSchemeStdio,
			internal.CompressionTypeNone,
		),
		"-",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatJSON,
			"",
			internal.FileSchemeStdio,
			internal.CompressionTypeNone,
		),
		"-#format=json",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"",
			internal.FileSchemeNull,
			internal.CompressionTypeNone,
		),
		app.DevNullFilePath,
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"path/to/dir",
			internal.FileSchemeLocal,
			internal.CompressionTypeNone,
		),
		"path/to/dir#format=bin",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"path/to/dir",
			internal.FileSchemeLocal,
			internal.CompressionTypeNone,
		),
		"path/to/dir#format=bin,compression=none",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"path/to/dir",
			internal.FileSchemeLocal,
			internal.CompressionTypeGzip,
		),
		"path/to/dir#format=bin,compression=gzip",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"/path/to/dir",
			internal.GitSchemeLocal,
			git.NewBranchName("main"),
			false,
			1,
			"",
		),
		"/path/to/dir#branch=main,format=git",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"/path/to/dir",
			internal.GitSchemeLocal,
			git.NewBranchName("main/foo"),
			false,
			1,
			"",
		),
		"/path/to/dir#format=git,branch=main/foo",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir",
			internal.GitSchemeLocal,
			git.NewTagName("main/foo"),
			false,
			1,
			"",
		),
		"path/to/dir#tag=main/foo,format=git",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir",
			internal.GitSchemeLocal,
			git.NewTagName("main/foo"),
			false,
			1,
			"",
		),
		"path/to/dir#format=git,tag=main/foo",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir",
			internal.GitSchemeLocal,
			git.NewTagName("main/foo"),
			true,
			1,
			"",
		),
		"path/to/dir#format=git,tag=main/foo,recurse_submodules=true",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir",
			internal.GitSchemeLocal,
			git.NewTagName("main/foo"),
			false,
			1,
			"",
		),
		"path/to/dir#format=git,tag=main/foo,recurse_submodules=false",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir",
			internal.GitSchemeLocal,
			git.NewRefName("refs/remotes/origin/HEAD"),
			false,
			50,
			"",
		),
		"path/to/dir#format=git,ref=refs/remotes/origin/HEAD",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedGitRef(
			formatGit,
			"path/to/dir",
			internal.GitSchemeLocal,
			git.NewRefName("refs/remotes/origin/HEAD"),
			false,
			10,
			"",
		),
		"path/to/dir#format=git,ref=refs/remotes/origin/HEAD,depth=10",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTargz,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeGzip,
			1,
			"",
		),
		"path/to/file#format=targz,strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeNone,
			1,
			"",
		),
		"path/to/file#format=tar,strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeNone,
			1,
			"",
		),
		"path/to/file#format=tar,strip_components=1,compression=none",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeGzip,
			1,
			"",
		),
		"path/to/file#format=tar,strip_components=1,compression=gzip",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatZip,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.ArchiveTypeZip,
			internal.CompressionTypeNone,
			1,
			"",
		),
		"path/to/file#format=zip,strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar.zst",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeZstd,
			0,
			"",
		),
		"path/to/file.tar.zst",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar.zst",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeZstd,
			1,
			"",
		),
		"path/to/file.tar.zst#strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatZip,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.ArchiveTypeZip,
			internal.CompressionTypeNone,
			1,
			"",
		),
		"path/to/file#format=zip,strip_components=1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file.tar.zst",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeZstd,
			0,
			"foo/bar",
		),
		"path/to/file.tar.zst#subdir=foo/bar",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedArchiveRef(
			formatTar,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.ArchiveTypeTar,
			internal.CompressionTypeZstd,
			1,
			"foo/bar",
		),
		"path/to/file#format=tar,strip_components=1,compression=zstd,subdir=foo/bar",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"path/to/file",
			internal.FileSchemeLocal,
			internal.CompressionTypeZstd,
		),
		"path/to/file#format=bin,compression=zstd",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"path/to/file.bin.zst",
			internal.FileSchemeLocal,
			internal.CompressionTypeZstd,
		),
		"path/to/file.bin.zst",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedModuleRef(
			formatMod,
			testNewModuleReference(
				t,
				"example.com",
				"foob",
				"bar",
				"v1",
			),
		),
		"example.com/foob/bar:v1",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedModuleRef(
			formatMod,
			testNewModuleReference(
				t,
				"example.com",
				"foob",
				"bar",
				bufmoduletesting.TestCommit,
			),
		),
		"example.com/foob/bar:"+bufmoduletesting.TestCommit,
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"github.com/path/to/file.bin",
			internal.FileSchemeHTTPS,
			internal.CompressionTypeNone,
		),
		"https://github.com/path/to/file.bin",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"github.com/path/to/file.ext",
			internal.FileSchemeHTTPS,
			internal.CompressionTypeNone,
		),
		"https://github.com/path/to/file.ext#format=bin",
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBin,
			"gitlab.com/api/v4/projects/foo/packages/generic/proto/0.0.1/proto.bin?private_token=bar",
			internal.FileSchemeHTTPS,
			internal.CompressionTypeNone,
		),
		"https://gitlab.com/api/v4/projects/foo/packages/generic/proto/0.0.1/proto.bin?private_token=bar#format=bin",
	)
}

func TestGetParsedRefError(t *testing.T) {
	testGetParsedRefError(
		t,
		internal.NewInvalidPathError(formatDir, "-"),
		"-#format=dir",
	)
	testGetParsedRefError(
		t,
		internal.NewInvalidPathError(formatGit, "-"),
		"-#format=git,branch=main",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyGitBranchAndTagError(),
		"path/to/foo#format=git,branch=foo,tag=bar",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyGitBranchAndTagError(),
		"path/to/foo#format=git,branch=foo,tag=bar,ref=baz",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyTagWithRefError(),
		"path/to/foo#format=git,tag=foo,ref=bar",
	)
	testGetParsedRefError(
		t,
		internal.NewDepthParseError("bar"),
		"path/to/foo#format=git,depth=bar",
	)
	testGetParsedRefError(
		t,
		internal.NewDepthZeroError(),
		"path/to/foo#format=git,ref=foor,depth=0",
	)
	testGetParsedRefError(
		t,
		internal.NewPathUnknownGzError("path/to/foo.gz"),
		"path/to/foo.gz",
	)
	testGetParsedRefError(
		t,
		internal.NewPathUnknownGzError("path/to/foo.bar.gz"),
		"path/to/foo.bar.gz",
	)
	testGetParsedRefError(
		t,
		internal.NewFormatOverrideNotAllowedForDevNullError(app.DevNullFilePath),
		fmt.Sprintf("%s#format=bin", app.DevNullFilePath),
	)
	testGetParsedRefError(
		t,
		internal.NewFormatUnknownError("bar"),
		"path/to/foo#format=bar",
	)
	testGetParsedRefError(
		t,
		internal.NewOptionsCouldNotParseStripComponentsError("foo"),
		"path/to/foo.tar.gz#strip_components=foo",
	)
	testGetParsedRefError(
		t,
		internal.NewCompressionUnknownError("foo"),
		"path/to/foo.tar.gz#compression=foo",
	)
	testGetParsedRefError(
		t,
		internal.NewOptionsInvalidKeyError("foo"),
		"path/to/foo.tar.gz#foo=bar",
	)
	testGetParsedRefError(
		t,
		internal.NewOptionsInvalidForFormatError(formatTar, "path/to/foo.tar.gz#branch=main"),
		"path/to/foo.tar.gz#branch=main",
	)
	testGetParsedRefError(
		t,
		internal.NewOptionsInvalidForFormatError(formatDir, "path/to/some/foo#strip_components=1"),
		"path/to/some/foo#strip_components=1",
	)
	testGetParsedRefError(
		t,
		internal.NewOptionsInvalidForFormatError(formatDir, "path/to/some/foo#compression=none"),
		"path/to/some/foo#compression=none",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyCompressionForZipError(),
		"path/to/foo.zip#compression=none",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyCompressionForZipError(),
		"path/to/foo.zip#compression=gzip",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyCompressionForZipError(),
		"path/to/foo#format=zip,compression=none",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyCompressionForZipError(),
		"path/to/foo#format=zip,compression=gzip",
	)
	testGetParsedRefError(
		t,
		internal.NewCannotSpecifyCompressionForZipError(),
		"path/to/foo#format=zip,compression=gzip",
	)
}

func testGetParsedRefSuccess(
	t *testing.T,
	expectedRef internal.ParsedRef,
	value string,
) {
	testGetParsedRef(
		t,
		expectedRef,
		nil,
		value,
	)
}

func testGetParsedRefError(
	t *testing.T,
	expectedErr error,
	value string,
) {
	testGetParsedRef(
		t,
		nil,
		expectedErr,
		value,
	)
}

func testGetParsedRef(
	t *testing.T,
	expectedParsedRef internal.ParsedRef,
	expectedErr error,
	value string,
) {
	t.Run(value, func(t *testing.T) {
		t.Parallel()
		parsedRef, err := newRefParser(zap.NewNop()).getParsedRef(
			context.Background(),
			value,
			allFormats,
		)
		if expectedErr != nil {
			if err == nil {
				assert.Equal(t, nil, parsedRef, "expected error")
			} else {
				assert.Equal(t, expectedErr, err)
			}
		} else {
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, expectedParsedRef, parsedRef)
			}
		}
	})
}

func testNewModuleReference(
	t *testing.T,
	remote string,
	owner string,
	module string,
	reference string,
) bufmoduleref.ModuleReference {
	moduleReference, err := bufmoduleref.NewModuleReference(remote, owner, module, reference)
	require.NoError(t, err)
	return moduleReference
}
