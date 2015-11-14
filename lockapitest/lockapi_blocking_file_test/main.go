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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/giancosta86/LockAPI/lockapi"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Argument required: <lock file>")
	}
	lockFilePath := os.Args[1]

	lockFile, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer lockFile.Close()

	err = lockapi.LockFile(lockFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press ENTER to unlock...")
	fmt.Scanln()

	err = lockapi.UnlockFile(lockFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press ENTER to end...")
	fmt.Scanln()
}
