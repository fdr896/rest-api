package bolt

import (
	"sort"
	"testing"
	"time"

	"github.com/fdr896/code_storage/core"
)

var (
	bucketName = []byte("tmp")
	path       = "../tmp/TemporaryStorageBolt.db"
	randomID   = "random ID"
	tmpCode    = &core.Code{
		ID:       "test id",
		Source:   "test code",
		Language: "test lang",
		Date:     time.Now(),
	}
	tmpCodeList = []*core.Code{
		&core.Code{
			ID:       "second test id",
			Source:   "second test code",
			Language: "second test lang",
		},
		&core.Code{
			ID:       "test id",
			Source:   "test code",
			Language: "test lang",
		},
	}
)

func failTest(t *testing.T, err error) {
	t.Errorf("test failed because: %v", err)
}

func TestNewCodeStorage(t *testing.T) {
	cs, err := NewCodeStorage(bucketName, path)
	defer cs.Close()
	if err != nil {
		failTest(t, err)
	}
}

func TestAdd(t *testing.T) {
	cs, _ := NewCodeStorage(bucketName, path)
	defer cs.Close()

	if err := cs.Add(tmpCode); err != nil {
		failTest(t, err)
	}
}

func TestGetWithExistingCode(t *testing.T) {
	cs, _ := NewCodeStorage(bucketName, path)
	defer cs.Close()

	cs.Add(tmpCode)

	if _, err := cs.Get(tmpCode.ID); err != nil {
		failTest(t, err)
	}
}

func TestGetWithNonExistingCode(t *testing.T) {
	cs, _ := NewCodeStorage(bucketName, path)
	defer cs.Close()

	if _, err := cs.Get(randomID); err != core.ErrNotFound {
		t.Errorf("database didn't recognise not existing code")
	}
}

func TestGetAllWithEmptyList(t *testing.T) {
	cs, _ := NewCodeStorage(bucketName, path)
	defer cs.Close()

	if _, err := cs.GetAll(); err != nil {
		t.Errorf("database failed when it was empty and error occured: %v", err)
	}
}

func TestGetAll(t *testing.T) {
	cs, _ := NewCodeStorage(bucketName, path)
	defer cs.Close()

	for _, code := range tmpCodeList {
		cs.Add(code)
	}

	receivedCodeList, _ := cs.GetAll()

	if len(receivedCodeList) != len(tmpCodeList) {
		t.Errorf("received and template lists are not equal because of different length: %d and %d", len(receivedCodeList), len(tmpCodeList))
	}

	sort.Slice(receivedCodeList, func(i, j int) bool {
		if receivedCodeList[i].ID != receivedCodeList[j].ID {
			return receivedCodeList[i].ID < receivedCodeList[j].ID
		}
		if receivedCodeList[i].Source != receivedCodeList[j].Source {
			return receivedCodeList[i].Source < receivedCodeList[j].Source
		}
		return receivedCodeList[i].Language < receivedCodeList[j].Language
	})

	for idx := range tmpCodeList {
		if tmpCodeList[idx].Source != receivedCodeList[idx].Source {
			t.Errorf("received and template lists are not equal")
		}

		if tmpCodeList[idx].Language != receivedCodeList[idx].Language {
			t.Errorf("received and template lists are not equal")
		}

		if tmpCodeList[idx].ID != receivedCodeList[idx].ID {
			t.Errorf("received and template lists are not equal")
		}
	}
}

func TestDeleteExistingCode(t *testing.T) {
	cs, _ := NewCodeStorage(bucketName, path)
	defer cs.Close()

	cs.Add(tmpCode)

	if err := cs.Delete(tmpCode.ID); err != nil {
		t.Errorf("failed to delte existing code because: %v", err)
	}
}

func TestDeleteNotExistingCode(t *testing.T) {
	cs, _ := NewCodeStorage(bucketName, path)
	defer cs.Close()

	if err := cs.Delete(randomID); err != nil {
		t.Errorf("failed to handle deleting not existing code because: %v", err)
	}
}
