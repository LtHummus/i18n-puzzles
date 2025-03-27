package input

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Kind int

const (
	RealInput Kind = iota
	TestInput
)

func (k Kind) String() string {
	if k == RealInput {
		return "REAL"
	} else {
		return "TEST"
	}
}

const (
	DirectoryName = ".i18n-puzzles"
	TokenFileName = ".token"

	BaseURL = "https://i18n-puzzles.com"

	cookieName = "sessionid"
)

func readToken(directory string) (string, error) {
	tokenFile := filepath.Join(directory, TokenFileName)

	f, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("downloader: readToken: could not read token file: %s: %w", tokenFile, err)
	}

	return strings.TrimSpace(string(f)), nil
}

func getPuzzleDirectory() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("downloader: getPuzzleDirectory: could not get home directory: %w", err)
	}

	puzzleDirectory := filepath.Join(dir, DirectoryName)

	_, err = os.Stat(puzzleDirectory)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(puzzleDirectory, 0755)
		if err != nil {
			return "", fmt.Errorf("downloader: getPuzzleDirectory: could not create home directory: %w", err)
		}
		fmt.Fprintf(os.Stderr, "created puzzle directory: %s\n", puzzleDirectory)
	}

	return puzzleDirectory, nil
}

func getInputFile(directory string, num int, kind Kind) string {
	var filename string
	if kind == TestInput {
		filename = fmt.Sprintf("%02d-test.txt", num)
	} else {
		filename = fmt.Sprintf("%02d.txt", num)
	}

	return filepath.Join(directory, filename)
}

func downloadInput(ctx context.Context, num int, kind Kind, token string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "downloading input for puzzle %d (input kind = %s)\n", num, kind)
	var remoteInputURL string
	if kind == TestInput {
		remoteInputURL = fmt.Sprintf("%s/puzzle/%d/test-input", BaseURL, num)
	} else {
		remoteInputURL = fmt.Sprintf("%s/puzzle/%d/input", BaseURL, num)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, remoteInputURL, nil)
	if err != nil {
		return nil, fmt.Errorf("downloader: downloadInput: could not build request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: token,
	})

	req.Header.Set("User-Agent", "i18n-puzzle downloader by LtHummus <lthummus.com>")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloader: downloadInput: could not make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	inputBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("downloader: downloadInput: could not read HTTP response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloader: downloadInput: non-200 response from server: %s", string(inputBytes))
	}

	return inputBytes, nil
}

func getPuzzleInput(ctx context.Context, num int, kind Kind, directory string, token string) ([]byte, error) {
	inputFile := getInputFile(directory, num, kind)
	_, err := os.Stat(inputFile)
	if err == nil {
		// !!! reverse of the normal if err test
		// this means the file exists and we can read it
		input, err := os.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("downloader: getPuzzleInput: could not read existing input file: %s: %w", inputFile, err)
		}
		return input, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		// if there's some error other than the file not existing, error out
		return nil, fmt.Errorf("downloader: getPuzzleInput: could not stat file: %s: %w", inputFile, err)
	}

	inputBytes, err := downloadInput(ctx, num, kind, token)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(inputFile, inputBytes, 0664)
	if err != nil {
		return nil, fmt.Errorf("downloader: getPuzzleInput: could not cache input: %s: %w", inputFile, err)
	}

	fmt.Fprintf(os.Stderr, "cached input at %s\n", inputFile)

	return inputBytes, nil
}

func GetInputBytes(ctx context.Context, num int, k Kind) ([]byte, error) {
	puzzleDir, err := getPuzzleDirectory()
	if err != nil {
		return nil, err
	}

	token, err := readToken(puzzleDir)
	if err != nil {
		return nil, err
	}

	return getPuzzleInput(ctx, num, k, puzzleDir, token)
}

func GetInputUTF8(ctx context.Context, num int, k Kind) (string, error) {
	puzzleDir, err := getPuzzleDirectory()
	if err != nil {
		return "", err
	}

	token, err := readToken(puzzleDir)
	if err != nil {
		return "", err
	}

	input, err := getPuzzleInput(ctx, num, k, puzzleDir, token)
	if err != nil {
		return "", err
	}

	return string(input), nil
}

func GetInputLinesUTF8(ctx context.Context, num int, k Kind) ([]string, error) {
	input, err := GetInputUTF8(ctx, num, k)
	if err != nil {
		return nil, err
	}

	potential := strings.Split(input, "\n")

	// last one might be empty
	if potential[len(potential)-1] == "" {
		return potential[:len(potential)-1], nil
	} else {
		return potential, nil
	}
}
