package network

import (
	"reflect"
	"testing"
)

func Test_parseIfconfig(t *testing.T) {
	var ifconfigOut = []byte{0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x30, 0x3a, 0x20, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x3d, 0x38, 0x38, 0x36, 0x33, 0x3c, 0x55, 0x50, 0x2c, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x2c, 0x53, 0x4d, 0x41, 0x52, 0x54, 0x2c, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x2c, 0x53, 0x49, 0x4d, 0x50, 0x4c, 0x45, 0x58, 0x2c, 0x4d, 0x55, 0x4c, 0x54, 0x49, 0x43, 0x41, 0x53, 0x54, 0x3e, 0x20, 0x6d, 0x74, 0x75, 0x20, 0x31, 0x35, 0x30, 0x30, 0xa, 0x9, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3d, 0x36, 0x33, 0x3c, 0x52, 0x58, 0x43, 0x53, 0x55, 0x4d, 0x2c, 0x54, 0x58, 0x43, 0x53, 0x55, 0x4d, 0x2c, 0x54, 0x53, 0x4f, 0x34, 0x2c, 0x54, 0x53, 0x4f, 0x36, 0x3e, 0xa, 0x9, 0x65, 0x74, 0x68, 0x65, 0x72, 0x20, 0x63, 0x65, 0x3a, 0x30, 0x30, 0x3a, 0x34, 0x34, 0x3a, 0x36, 0x30, 0x3a, 0x35, 0x39, 0x3a, 0x30, 0x35, 0x20, 0xa, 0x9, 0x69, 0x6e, 0x65, 0x74, 0x20, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x20, 0x6e, 0x65, 0x74, 0x6d, 0x61, 0x73, 0x6b, 0x20, 0x30, 0x78, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x30, 0x30, 0x20, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x20, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x32, 0x35, 0x35, 0xa, 0x9, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0xa, 0x9, 0x9, 0x69, 0x64, 0x20, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x20, 0x30, 0x20, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x30, 0x20, 0x66, 0x77, 0x64, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x20, 0x30, 0xa, 0x9, 0x9, 0x6d, 0x61, 0x78, 0x61, 0x67, 0x65, 0x20, 0x30, 0x20, 0x68, 0x6f, 0x6c, 0x64, 0x63, 0x6e, 0x74, 0x20, 0x30, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x20, 0x73, 0x74, 0x70, 0x20, 0x6d, 0x61, 0x78, 0x61, 0x64, 0x64, 0x72, 0x20, 0x31, 0x30, 0x30, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x20, 0x31, 0x32, 0x30, 0x30, 0xa, 0x9, 0x9, 0x72, 0x6f, 0x6f, 0x74, 0x20, 0x69, 0x64, 0x20, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x3a, 0x30, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x20, 0x30, 0x20, 0x69, 0x66, 0x63, 0x6f, 0x73, 0x74, 0x20, 0x30, 0x20, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x30, 0xa, 0x9, 0x9, 0x69, 0x70, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x20, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x20, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x20, 0x30, 0x78, 0x32, 0xa, 0x9, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x3a, 0x20, 0x65, 0x6e, 0x33, 0x20, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x3d, 0x33, 0x3c, 0x4c, 0x45, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x2c, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x56, 0x45, 0x52, 0x3e, 0xa, 0x9, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x69, 0x66, 0x6d, 0x61, 0x78, 0x61, 0x64, 0x64, 0x72, 0x20, 0x30, 0x20, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x31, 0x37, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x20, 0x30, 0x20, 0x70, 0x61, 0x74, 0x68, 0x20, 0x63, 0x6f, 0x73, 0x74, 0x20, 0x30, 0xa, 0x9, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x3a, 0x20, 0x65, 0x6e, 0x31, 0x20, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x3d, 0x33, 0x3c, 0x4c, 0x45, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x2c, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x56, 0x45, 0x52, 0x3e, 0xa, 0x9, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x69, 0x66, 0x6d, 0x61, 0x78, 0x61, 0x64, 0x64, 0x72, 0x20, 0x30, 0x20, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x31, 0x35, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x20, 0x30, 0x20, 0x70, 0x61, 0x74, 0x68, 0x20, 0x63, 0x6f, 0x73, 0x74, 0x20, 0x30, 0xa, 0x9, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x3a, 0x20, 0x65, 0x6e, 0x32, 0x20, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x3d, 0x33, 0x3c, 0x4c, 0x45, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x2c, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x56, 0x45, 0x52, 0x3e, 0xa, 0x9, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x69, 0x66, 0x6d, 0x61, 0x78, 0x61, 0x64, 0x64, 0x72, 0x20, 0x30, 0x20, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x31, 0x36, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x20, 0x30, 0x20, 0x70, 0x61, 0x74, 0x68, 0x20, 0x63, 0x6f, 0x73, 0x74, 0x20, 0x30, 0xa, 0x9, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x3a, 0x20, 0x65, 0x6e, 0x34, 0x20, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x3d, 0x33, 0x3c, 0x4c, 0x45, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x2c, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x56, 0x45, 0x52, 0x3e, 0xa, 0x9, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x69, 0x66, 0x6d, 0x61, 0x78, 0x61, 0x64, 0x64, 0x72, 0x20, 0x30, 0x20, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x31, 0x34, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x20, 0x30, 0x20, 0x70, 0x61, 0x74, 0x68, 0x20, 0x63, 0x6f, 0x73, 0x74, 0x20, 0x30, 0xa, 0x9, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x20, 0x63, 0x61, 0x63, 0x68, 0x65, 0x3a, 0xa, 0x9, 0x6e, 0x64, 0x36, 0x20, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3d, 0x32, 0x30, 0x31, 0x3c, 0x50, 0x45, 0x52, 0x46, 0x4f, 0x52, 0x4d, 0x4e, 0x55, 0x44, 0x2c, 0x44, 0x41, 0x44, 0x3e, 0xa, 0x9, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x3a, 0x20, 0x3c, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x20, 0x74, 0x79, 0x70, 0x65, 0x3e, 0xa, 0x9, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x3a, 0x20, 0x69, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0xa}

	type args struct {
		ifconfig string
	}

	tests := []struct {
		name string
		args args
		want *Bridge
	}{
		{
			name: "bridge0",
			args: args{ifconfig: string(ifconfigOut)},
			want: &Bridge{
				Members: []string{"en3", "en1", "en2", "en4"},
				IP:      "10.0.0.1",
				Netmask: []byte{255, 255, 255, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseIfconfig(tt.args.ifconfig)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIfconfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sliceDiff(t *testing.T) {
	type args struct {
		x []string
		y []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want2 []string
	}{
		{
			name: "No difference with duplicates",
			args: args{x: []string{"a", "a", "c"}, y: []string{"c", "a", "c"}},
			want: nil,
		},
		{
			name:  "Equal not in the same order",
			args:  args{x: []string{"z", "z", "x"}, y: []string{"x", "z", "z"}},
			want:  nil,
			want2: nil,
		},
		{
			name:  "One new element",
			args:  args{x: []string{"x", "z"}, y: []string{"z"}},
			want:  []string{"x"},
			want2: nil,
		},
		{
			name:  "Two new elements",
			args:  args{x: []string{"x", "z"}, y: []string{}},
			want:  []string{"x", "z"},
			want2: nil,
		},
		{
			name:  "No new elements, two existing",
			args:  args{x: []string{}, y: []string{"x", "z"}},
			want:  nil,
			want2: []string{"x", "z"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := sliceDiff(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceDiff() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("sliceDiff() = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}

// Test wraps exec.Command to redirect it to calling the TestIfconfigHelperProcess below.
//type Test struct{}
//
//func (t Test) Output(command string, args ...string) ([]byte, error) {
//	cs := []string{"-test.run=TestIfconfigHelperProcess", "--"}
//	cs = append(cs, args...)
//	cmd := exec.Command(os.Args[0], cs...)
//	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
//	out, err := cmd.Output()
//	return out, err
//}
//
//func (t Test) Run(command string, args ...string) error {
//	cs := []string{"-test.run=TestIfconfigHelperProcess", "--"}
//	cs = append(cs, args...)
//	cmd := exec.Command(os.Args[0], cs...)
//	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
//	err := cmd.Run()
//	return err
//}
//
//func (t Test) CombinedOutput(command string, args ...string) ([]byte, error) {
//	cs := []string{"-test.run=TestIfconfigHelperProcess", "--"}
//	cs = append(cs, args...)
//	cmd := exec.Command(os.Args[0], cs...)
//	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
//	out, err := cmd.CombinedOutput()
//	return out, err
//}
//
//// TestHelperProcess isn't a real test. It's used as a helper process
//func TestIfconfigHelperProcess(*testing.T) {
//	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
//		return
//	}
//	defer os.Exit(0)
//
//	args := os.Args
//	for len(args) > 0 {
//		if args[0] == "--" {
//			args = args[1:]
//			break
//		}
//		args = args[1:]
//	}
//	if len(args) == 0 {
//		fmt.Fprintf(os.Stderr, "No command\n")
//		os.Exit(2)
//	}
//
//	cmd, args := args[0], args[1:]
//	if cmd != "ifconfig" {
//		fmt.Fprintf(os.Stderr, "Unknown command\n")
//		os.Exit(2)
//	}
//	if args[1] != "bridge0" {
//		fmt.Fprintf(os.Stderr, "ifconfig: interface %s does not exist\n", args[1])
//		os.Exit(1)
//	}
//
//	if len(args) == 1 {
//		fmt.Printf("%s", ifconfigOut)
//	}
//}
