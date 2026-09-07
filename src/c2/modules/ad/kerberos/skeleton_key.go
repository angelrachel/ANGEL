package kerberos

import (
"crypto/md5"
"encoding/hex"
)

type SkeletonKey struct {
Domain   string
}

func NewSkeletonKey(domain string) *SkeletonKey {
return &SkeletonKey{
Domain: domain,
}
}

func (s *SkeletonKey) CreateSkeletonKey(password string) string {
// Skeleton Key is a master password that works for all users
// This is a simplified implementation
hash := md5.Sum([]byte(password))
return hex.EncodeToString(hash[:])
}

func (s *SkeletonKey) GetSkeletonHash() string {
// Mimikatz skeleton key uses "mimikatz" as the password by default
return s.CreateSkeletonKey("mimikatz")
}
