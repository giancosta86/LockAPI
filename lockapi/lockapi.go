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

/*
Minimalist Go API dedicated to cross-platform locking based on OS-specific
system calls.
*/
package lockapi

import "os"

/*
Exclusively locks an existing file.
Returns nil if the lock could be obtained successfully,
otherwise returns an error immediately.
*/
func TryLockFile(file *os.File) (err error) {
	return tryLockFileImpl(file)
}

/*
Exclusively locks an existing file, blocking until the lock can be obtained or
an error occurs.
Returns nil on success, an error on failure.
*/
func LockFile(file *os.File) (err error) {
	return lockFileImpl(file)
}

/*
Releases the lock on the given file.
Returns nil on success, an error on failure.
*/
func UnlockFile(file *os.File) (err error) {
	return unlockFileImpl(file)
}
