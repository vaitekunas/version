package main

import (
  "testing"
)

func TestLarger(t *testing.T) {

  tests := []struct{
    v *Version
    w *Version
    larger bool
  }{
    {&Version{Major:0,Minor:0,Patch:1},&Version{Major:0,Minor:0,Patch:0},true},
    {&Version{Major:0,Minor:0,Patch:0},&Version{Major:0,Minor:0,Patch:1},false},
    {&Version{Major:0,Minor:1,Patch:0},&Version{Major:0,Minor:0,Patch:10},true},
    {&Version{Major:10,Minor:1,Patch:3},&Version{Major:0,Minor:10,Patch:10},true},
    {&Version{Major:10,Minor:100,Patch:300},&Version{Major:20,Minor:10,Patch:10},false},
    {&Version{Major:1,Minor:0,Patch:0},&Version{Major:1,Minor:0,Patch:0, Special: "rc1"},true},
    {&Version{Major:1,Minor:0,Patch:0},&Version{Major:1,Minor:0,Patch:0, Special: "alpha.rc1"},true},
    {&Version{Major:1,Minor:0,Patch:0,Special:"alpha.rc2"},&Version{Major:1,Minor:0,Patch:0, Special: "beta.rc1"},false},
    {&Version{Major:1,Minor:0,Patch:0,Special:"alpha"},&Version{Major:1,Minor:0,Patch:0, Special: "beta.rc1"},false},
    {&Version{Major:1,Minor:0,Patch:0,Special:"beta"},&Version{Major:1,Minor:0,Patch:0, Special: "beta.rc1"},false},
    {&Version{Major:1,Minor:0,Patch:0,Special:"gamma"},&Version{Major:1,Minor:0,Patch:0, Special: "beta.rc1"},true},
    {&Version{Major:1,Minor:0,Patch:0,Special:"gamma.rc2"},&Version{Major:1,Minor:0,Patch:0, Special: "beta.rc3"},true},
    {&Version{Major:1,Minor:0,Patch:0,Special:"beta.1.rc2"},&Version{Major:1,Minor:0,Patch:0, Special: "beta.rc1.1"},false},
    {&Version{Major:1,Minor:0,Patch:0},&Version{Major:1,Minor:0,Patch:1, Special: "alpha.rc1"}, false},
  }

  for i, test := range tests {
    if Larger(test.v,test.w) != test.larger {
      t.Errorf("TestLarger: test %d failed",i+1)
    }
  }

}
