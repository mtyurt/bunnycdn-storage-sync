package syncer

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mtyurt/bunnycdn-storage-sync/api"
)

// BCDNSyncer is the service that runs sync operation
type BCDNSyncer struct {
	API    api.BCDNStorage
	DryRun bool
}

// Sync synchronizes sourcePath with storage zone efficiently
func (s *BCDNSyncer) Sync(sourcePath string) error {
	objMap := make(map[string]api.BCDNObject)

	metrics := struct {
		total        int
		newFile      int
		modifiedFile int
		deletedFile  int
	}{0, 0, 0, 0}
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		log.Printf("DEBUG: checking path: %q\n", path)
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return s.fetchDirectory(objMap, relPath)
		}
		metrics.total += 1
		fileContent, fsChecksum, err := getFileContent(path)

		obj, ok := objMap[relPath]
		if !ok {
			log.Printf("DEBUG: %s not found in storage, uploading...\n", relPath)
			metrics.newFile += 1
			return s.uploadFile(relPath, fileContent, fsChecksum)
		}

		fsChecksum = strings.ToLower(fsChecksum)
		objChecksum := strings.ToLower(obj.Checksum)
		if fsChecksum != objChecksum {
			metrics.modifiedFile += 1
			log.Printf("DEBUG: %s checksum doesn't match, local: [%s] remote: [%s]\n", relPath, fsChecksum, objChecksum)
			delete(objMap, relPath)
			return s.uploadFile(relPath, fileContent, fsChecksum)
		}
		log.Printf("DEBUG: %s matches with remote storage, skipping.\n", relPath)

		delete(objMap, relPath)

		return nil
	})
	for relPath, obj := range objMap {
		if !obj.IsDirectory {
			metrics.deletedFile += 1
			log.Printf("%s object must be deleted.\n", relPath)
			s.deletePath(relPath)
		}
	}

	log.Printf("Total files: %d New files: %d Modified files: %d Deleted files: %d\n", metrics.total, metrics.newFile, metrics.modifiedFile, metrics.deletedFile)

	return err
}

// fetchDirectory fetches path from BunnyCDN API & stores all objects in a map with their path as key
func (s *BCDNSyncer) fetchDirectory(objMap map[string]api.BCDNObject, path string) error {
	log.Printf("DEBUG: Fetching directory %s\n", path)
	objects, err := s.API.List(path)
	if err != nil {
		return err
	}
	zoneName := s.API.ZoneName
	for _, obj := range objects {
		objPath := strings.TrimPrefix(obj.Path+obj.ObjectName, "/"+zoneName+"/")
		objMap[objPath] = obj
		// log.Printf("TRACE: Mapping %s to %v\n", objPath, obj)
	}
	return nil
}

func (s *BCDNSyncer) uploadFile(path string, content []byte, checksum string) error {
	log.Printf("Uploading file %s with checksum %s\n", path, checksum)
	if s.DryRun {
		return nil
	}
	return s.API.Upload(path, content, checksum)
}

func (s *BCDNSyncer) deletePath(path string) error {
	log.Printf("Deleting file %s\n", path)
	if s.DryRun {
		return nil
	}
	return s.API.Delete(path)
}

func getFileContent(path string) ([]byte, string, error) {

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	checksum := sha256.Sum256(fileContent)
	return fileContent, fmt.Sprintf("%x", checksum), nil
}
