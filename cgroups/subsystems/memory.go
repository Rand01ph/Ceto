package subsystems


type MemorySubSystem struct {

}

// 设置cgroupPath对应的cgroup的内存限制
func (s *MemorySubSystem) Set (cgroupPath string, res *ResourceConfig) error {

	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {

	}
}