/*§
  ===========================================================================
  LockAPI
  ===========================================================================
  Copyright (C) 2015 Gianluca Costa
  ===========================================================================
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  ===========================================================================
*/

package lockapitest

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/giancosta86/LockAPI/lockapi"
	"github.com/giancosta86/LockAPI/testutils"
)

func TestNonBlockingFileLockingCalls(t *testing.T) {
	lockFile, err := ioutil.TempFile(os.TempDir(), "LockAPI_TestFile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		lockFile.Close()
		os.Remove(lockFile.Name())
	}()

	err = lockapi.TryLockFile(lockFile)
	if err != nil {
		t.Fatal(err)
	}

	err = lockapi.UnlockFile(lockFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBlockingFileLockingCalls(t *testing.T) {
	lockFile, err := ioutil.TempFile(os.TempDir(), "LockAPI_TestFile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		lockFile.Close()
		os.Remove(lockFile.Name())
	}()

	err = lockapi.LockFile(lockFile)
	if err != nil {
		t.Fatal(err)
	}

	err = lockapi.UnlockFile(lockFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNonBlockingFileLocks(t *testing.T) {
	testProgramPath, err := testutils.GoInstall("github.com/giancosta86/LockAPI/lockapitest/lockapi_non_blocking_file_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(testProgramPath)

	lockFile, err := ioutil.TempFile(os.TempDir(), "LockAPI_TestFile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		lockFile.Close()
		os.Remove(lockFile.Name())
	}()

	firstProcess := exec.Command(testProgramPath, lockFile.Name())

	firstInput, err := firstProcess.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}

	firstOutput, err := firstProcess.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	err = firstProcess.Start()
	if err != nil {
		t.Fatal(err)
	}

	//Reading the output will pause this process until the first process
	//has actually acquired the lock
	buffer := make([]byte, 5)
	outBytesCount, err := firstOutput.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if outBytesCount <= 0 {
		t.Fatal("Error while reading the output")
	}

	secondProcess := exec.Command(testProgramPath, lockFile.Name())
	err = secondProcess.Run()
	if err != nil {
		if _, isExitError := err.(*exec.ExitError); isExitError {
			//Just do nothing: the second process MUST fail
		} else {
			t.Fatal(err)
		}
	}

	if secondProcess.ProcessState.Success() {
		log.Fatal("The second process should fail")
	}

	//Now unlocking the file...
	_, err = firstInput.Write([]byte("\n"))
	if err != nil {
		t.Fatal(err)
	}

	thirdProcess := exec.Command(testProgramPath, lockFile.Name())
	err = thirdProcess.Run()
	if err != nil {
		t.Fatal(err)
	}

	if !thirdProcess.ProcessState.Success() {
		log.Fatal("The third process should succeed")
	}

	//Now terminating the first process...
	_, err = firstInput.Write([]byte("\n"))
	if err != nil {
		t.Fatal(err)
	}

	err = firstProcess.Wait()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBlockingFileLocks(t *testing.T) {
	testProgramPath, err := testutils.GoInstall("github.com/giancosta86/LockAPI/lockapitest/lockapi_blocking_file_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(testProgramPath)

	lockFile, err := ioutil.TempFile(os.TempDir(), "LockAPI_TestFile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		lockFile.Close()
		os.Remove(lockFile.Name())
	}()

	firstProcess := exec.Command(testProgramPath, lockFile.Name())

	firstInput, err := firstProcess.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}

	firstOutput, err := firstProcess.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	err = firstProcess.Start()
	if err != nil {
		t.Fatal(err)
	}

	//Reading the output will pause this process until the first process
	//has actually acquired the lock
	buffer := make([]byte, 5)
	firstOutputBytesCount, err := firstOutput.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if firstOutputBytesCount <= 0 {
		t.Fatal("Error while reading the first output")
	}

	secondProcess := exec.Command(testProgramPath, lockFile.Name())

	secondOutput, err := secondProcess.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = secondProcess.Start()
	if err != nil {
		t.Fatal(err)
	}

	firstSecondChannel := make(chan interface{})
	defer close(firstSecondChannel)

	go func() {
		secondOutputBytesCount, err := secondOutput.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		if secondOutputBytesCount <= 0 {
			t.Fatal("Error while reading the second output")
		}

		firstSecondChannel <- nil
	}()

	select {
	case <-firstSecondChannel:
		log.Fatal("The second process should NOT print output before succeeding in acquiring the lock")
	case <-time.After(3 * time.Second):
		//Just do nothing
	}

	//Now unlocking the file...
	_, err = firstInput.Write([]byte("\n"))
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-firstSecondChannel:
		//Just do nothing
	case <-time.After(5 * time.Second):
		log.Fatal("The second process should be notified of the released lock")
	}

	err = secondProcess.Wait()
	if err != nil {
		t.Fatal(err)
	}

	if !secondProcess.ProcessState.Success() {
		log.Fatal("The second process should succeed")
	}

	thirdProcess := exec.Command(testProgramPath, lockFile.Name())
	err = thirdProcess.Run()
	if err != nil {
		t.Fatal(err)
	}

	if !thirdProcess.ProcessState.Success() {
		log.Fatal("The third process should succeed")
	}

	//Now terminating the first process...
	_, err = firstInput.Write([]byte("\n"))
	if err != nil {
		t.Fatal(err)
	}

	err = firstProcess.Wait()
	if err != nil {
		t.Fatal(err)
	}
}
