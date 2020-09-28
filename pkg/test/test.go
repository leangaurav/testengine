package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

type ProgrammingLang int
const (
	PYTHON ProgrammingLang  = iota
)

type Test struct {
	cmd *exec.Cmd
	codeFileName string

	// ExecutionDir should be an absolute path
	ID string
	Code string
	Language ProgrammingLang
	ExecutionDir string
	TestCases  []TestCase
	
	maxResource Resources
}

type Resources struct {
	Time time.Duration
}

func NewTest(id string, code string, language ProgrammingLang, executionDir string ,testCases []TestCase) *Test {
	return &Test {
		cmd:nil,

		ID: id,
		Code:code,
		Language: language,
		ExecutionDir: path.Join(executionDir, id),
		TestCases: testCases,

		maxResource: Resources {
			Time: time.Second * 2,
		},
	}
}

func(t *Test) Run() error {
	
	// save code to a temp file
	if err := t.init(); err != nil {
		return fmt.Errorf("error saving code: %w", err)
	}

	// execute code
	if err := t.execute(); err != nil {
		return fmt.Errorf("error executing code: %w", err)
	}

	// cleanup
	if err := t.cleanup(); err != nil {
		return fmt.Errorf("error performing cleanup: %w", err)
	}
 
	return nil
}

func(t *Test) init() error {
	var (
		err error
	)

	if !path.IsAbs(t.ExecutionDir) {
		return errors.New(fmt.Sprintf("execution location'%s' is not absolute path", t.ExecutionDir))
	}

	t.codeFileName, err = getFileName(t.Language)
	if err != nil {
		return err
	}

	// create execution directory and save file
	err = os.Mkdir(t.ExecutionDir, 0700)
	if err != nil {
		return fmt.Errorf("unable to create execution dir: %w", err)
	}

	filePath := path.Join(t.ExecutionDir, t.codeFileName) 
	if err := t.saveCode(filePath); err != nil {
		return err
	}

	return nil
}

func (t *Test) saveCode(filePath string) error{
	var  (
		err error
	)
	
	file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to open code file: %w", err)
	} 

	_, err = file.WriteString(t.Code)
	if err != nil {
		return fmt.Errorf("Unable to write code to file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("Unable to close code file: %w", err)
	}
	return nil
}

func(t *Test) execute() error {
	var  (
		err error
	)

	cmd, args, err := t.getCommandParams()
	if err != nil {
		return err
	}

	// create command with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), t.maxResource.Time)
	defer cancel()
	t.cmd = exec.CommandContext(ctx, cmd, args...)
	
	fmt.Println("running command", cmd, args)
	// TODO: set permissions

	// Run the command and wait for execution
	if err := t.cmd.Run(); err != nil {
		return fmt.Errorf("error while executing command: %w", err)
	}

	return nil
}

func (t *Test) cleanup() error {
	// delete folder
	if err := os.RemoveAll(t.ExecutionDir); err != nil {
		return fmt.Errorf("unable to do cleanup: %w", err)
	}

	return nil
}

func(t *Test) getCommandParams() (string, []string, error) {
	baseCmd := ""

	switch t.Language {
	case PYTHON: baseCmd = "python"
	default: return "", []string{}, errors.New(fmt.Sprintf("no command mapping for language: %v", t.Language))
	}

	file := path.Join(t.ExecutionDir, t.codeFileName)
	args := []string{file}

	return baseCmd, args, nil
}

func getFileName(language ProgrammingLang) (string, error) {
	switch language {
	case PYTHON: return "main.py", nil
	}
	return "", errors.New(fmt.Sprintf("No file name mapping found for language: %v", language))
}
