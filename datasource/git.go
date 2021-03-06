package datasource

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitDataSource is the git data source object
type GitDataSource struct{}

func Push(from, to string) error {
	fmt.Printf("Pushing data from %s into %s...\n", from, to)
	//if err := createFolderIfNotExist(to); err != nil {
	//	return err
	//}
	if err := pushRepo(from,to); err != nil {
		return err
	}
	fmt.Print("Pushing complete.\n")
	return  nil
}

// Fetch creates the output folder, clears it and clones the repository there
func (ds *GitDataSource) Fetch(from, to string) ([]string, error) {
	fmt.Printf("Fetching data from %s into %s...\n", from, to)
	if err := createFolderIfNotExist(to); err != nil {
		return nil, err
	}
	if err := clearFolder(to); err != nil {
		return nil, err
	}
	if err := cloneRepo(to, from); err != nil {
		return nil, err
	}
	dirs, err := getContentFolders(to, ".md")
	if err != nil {
		return nil, err
	}
	fmt.Print("Fetching complete.\n")
	return dirs, nil
}

func pushRepo(path, repositoryURL string) error {
	cmdName := "git"
	initArgs := []string{"init", "."}
	cmd := exec.Command(cmdName, initArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error initializing git repository at %s: %v", path, err)
	}
	remoteArgs := []string{"remote", "add", "origin2", repositoryURL}
	cmd = exec.Command(cmdName, remoteArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting remote %s: %v", repositoryURL, err)
	}
	addArgs := []string{"add", "."}
	cmd = exec.Command(cmdName, addArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error add .: %v", err)
	}
	commitArgs := []string{"commit", "-m", "upadte"}
	cmd = exec.Command(cmdName, commitArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error commit : %v", err)
	}
	pullArgs := []string{"push", "origin2", "master", "--force"}
	cmd = exec.Command(cmdName, pullArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error pushing master at %s: %v", path, err)
	}
	return nil
}

func createFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(path, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory %s: %v", path, err)
			}
		} else {
			return fmt.Errorf("error accessing directory %s: %v", path, err)
		}
	}
	return nil
}

func clearFolder(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error accessing directory %s: %v", path, err)
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", path, err)
	}

	for _, name := range names {
		if err = os.RemoveAll(filepath.Join(path, name)); err != nil {
			return fmt.Errorf("error clearing file %s: %v", name, err)
		}
	}
	return nil
}

func cloneRepo(path, repositoryURL string) error {
	cmdName := "git"
	initArgs := []string{"init", "."}
	cmd := exec.Command(cmdName, initArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error initializing git repository at %s: %v", path, err)
	}
	remoteArgs := []string{"remote", "add", "origin", repositoryURL}
	cmd = exec.Command(cmdName, remoteArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting remote %s: %v", repositoryURL, err)
	}
	pullArgs := []string{"pull", "origin", "master"}
	cmd = exec.Command(cmdName, pullArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error pulling master at %s: %v", path, err)
	}
	return nil
}

func getContentFolders(path string, fileType string) ([]string, error) {
	var result []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), fileType) {
			result = append(result, path)
		}
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("file not exit %s: %v", path, err)
	}
	return result, nil
}
