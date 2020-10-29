/*===============================================================
*   Copyright (C) 2020 All rights reserved.
*
*   FileName：cpuidle_test.go
*   Author：WuGuoFu
*   Date： 2020-10-29
*   Description：
*
================================================================*/
package cpuidle

import (
	"sync"
	"testing"
)

func TestCpuStat(t *testing.T) {
	cp := new(CpuIdlePlugin)
	lock := new(sync.RWMutex)
	cp.lock = lock
	cp.updateCpuStat(1)
}
