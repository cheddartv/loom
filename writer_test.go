package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	main "github.com/cheddartv/loom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindStructIndexByPath", func() {
	input := []main.ParsedInput{
		{Path: "index.m3u8", AbsPath: "/root/index.m3u8", Include: true, Playlist: nil},
		{Path: "index2.m3u8", AbsPath: "/root/index2.m3u8", Include: true, Playlist: nil},
	}
	It("returns the index of the struct", func() {
		Expect(main.FindStructIndexByPath("/root/index.m3u8", input)).To(Equal(0))
	})
	It("returns -1 if the file is not present", func() {
		Expect(main.FindStructIndexByPath("not_present.m3u8", input)).To(Equal(-1))
	})
})

var _ = Describe("HandleEvent", func() {
	workingDir, _ := filepath.EvalSymlinks(os.Getenv("PWD"))
	input := []main.ParsedInput{
		{Path: "index.m3u8", AbsPath: workingDir + "/example/primary/index.m3u8", Include: false, Playlist: nil},
		{Path: "index.m3u8", AbsPath: workingDir + "/example/backup/index.m3u8", Include: true, Playlist: nil},
		{Path: "1.m3u8", AbsPath: workingDir + "/example/primary/1.m3u8", Include: false, Playlist: nil},
	}
	It("Updates the data struct on file creation", func() {
		change := main.Change{Path: "index.m3u8", AbsPath: workingDir + "/example/backup/index.m3u8", Type: "Create"}
		Expect(main.HandleEvent(change, input)[1].Include).To(BeTrue())
	})
	It("Updates the data struct on file removal", func() {
		change := main.Change{Path: "index.m3u8", AbsPath: workingDir + "/example/primary/index.m3u8", Type: "Remove"}
		Expect(main.HandleEvent(change, input)[0].Include).To(BeFalse())
	})
	It("non-tracked files do not change data struct", func() {
		change := main.Change{Path: "not_tracked.m3u8", AbsPath: workingDir + "/example/not_tracked.m3u8", Type: "Write"}
		Expect(main.HandleEvent(change, input)).Should(BeEquivalentTo(input))
	})
	It("Removes a playlist if an update makes it fail parsing", func() {
		change := main.Change{Path: "1.m3u8", AbsPath: workingDir + "/example/primary/1.m3u8", Type: "Write"}
		Expect(main.HandleEvent(change, input)[2].Include).To(BeFalse())
	})
	It("Updates the data struct for create on file not in the struct", func() {
		input := []main.ParsedInput{
			{Path: "index.m3u8", AbsPath: workingDir + "/example/primary/index.m3u8", Include: true, Playlist: nil},
			{Path: "1.m3u8", AbsPath: workingDir + "/example/primary/1.m3u8", Include: false, Playlist: nil},
		}
		change := main.Change{Path: "index.m3u8", AbsPath: workingDir + "/example/backup/index.m3u8", Type: "Create"}
		Expect(len(main.HandleEvent(change, input))).To(BeEquivalentTo(3))
	})
	It("Does not track a create on a playlist that fails to parse", func() {
		input := []main.ParsedInput{
			{Path: "index.m3u8", AbsPath: workingDir + "/example/primary/index.m3u8", Include: false, Playlist: nil},
		}
		change := main.Change{Path: "1.m3u8", AbsPath: workingDir + "/example/primary/1.m3u8", Type: "Create"}
		Expect(main.HandleEvent(change, input)).To(BeEquivalentTo(input))
	})
})

var _ = Describe("WriteManifest", func() {
	output := "tmp/index.m3u8"
	workingDir, _ := filepath.EvalSymlinks(os.Getenv("PWD"))
	mp1, _ := main.ImportPlaylist(workingDir + "/example/primary/index.m3u8")
	mp2, _ := main.ImportPlaylist(workingDir + "/example/backup/index.m3u8")
	manifests := []main.ParsedInput{
		{Path: "index.m3u8", AbsPath: workingDir + "/example/primary/index.m3u8", Include: true, Playlist: mp1},
		{Path: "index.m3u8", AbsPath: workingDir + "/example/backup/index.m3u8", Include: true, Playlist: mp2},
	}
	It("Generates and output file", func() {
		main.WriteManifest(manifests, main.CleanPath("tmp/index.m3u8"))
		Expect(output).Should(BeAnExistingFile())
	})

	It("sorts weaves and sorts by bitrate", func() {
		filebyteBuffer, _ := ioutil.ReadFile(output)
		filecontents := string(filebyteBuffer)
		Expect(filecontents).Should(MatchRegexp("#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=8171676,CODECS=\"avc1.4d4028,mp4a.40.5\",RESOLUTION=1920x1080\n../example/primary/1.m3u8\n#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=8171676,CODECS=\"avc1.4d4028,mp4a.40.5\",RESOLUTION=1920x1080\n../example/backup/1.m3u8\n#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=6332540,CODECS=\"avc1.4d401f,mp4a.40.5\",RESOLUTION=1280x720\n../example/primary/2.m3u8\n#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=6332540,CODECS=\"avc1.4d401f,mp4a.40.5\",RESOLUTION=1280x720\n../example/primary/3.m3u8\n#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=6332540,CODECS=\"avc1.4d401f,mp4a.40.5\",RESOLUTION=1280x720\n../example/backup/2.m3u8"))
	})

})
