// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.0
// source: proto/filetransfer.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FileChunk struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileName      string                 `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Data          []byte                 `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	ChunkNumber   int64                  `protobuf:"varint,3,opt,name=chunk_number,json=chunkNumber,proto3" json:"chunk_number,omitempty"`
	TotalChunks   int64                  `protobuf:"varint,4,opt,name=total_chunks,json=totalChunks,proto3" json:"total_chunks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileChunk) Reset() {
	*x = FileChunk{}
	mi := &file_proto_filetransfer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileChunk) ProtoMessage() {}

func (x *FileChunk) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filetransfer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileChunk.ProtoReflect.Descriptor instead.
func (*FileChunk) Descriptor() ([]byte, []int) {
	return file_proto_filetransfer_proto_rawDescGZIP(), []int{0}
}

func (x *FileChunk) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileChunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *FileChunk) GetChunkNumber() int64 {
	if x != nil {
		return x.ChunkNumber
	}
	return 0
}

func (x *FileChunk) GetTotalChunks() int64 {
	if x != nil {
		return x.TotalChunks
	}
	return 0
}

type FileTransferResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	BytesReceived int64                  `protobuf:"varint,3,opt,name=bytes_received,json=bytesReceived,proto3" json:"bytes_received,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileTransferResponse) Reset() {
	*x = FileTransferResponse{}
	mi := &file_proto_filetransfer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileTransferResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileTransferResponse) ProtoMessage() {}

func (x *FileTransferResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filetransfer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileTransferResponse.ProtoReflect.Descriptor instead.
func (*FileTransferResponse) Descriptor() ([]byte, []int) {
	return file_proto_filetransfer_proto_rawDescGZIP(), []int{1}
}

func (x *FileTransferResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *FileTransferResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *FileTransferResponse) GetBytesReceived() int64 {
	if x != nil {
		return x.BytesReceived
	}
	return 0
}

var File_proto_filetransfer_proto protoreflect.FileDescriptor

const file_proto_filetransfer_proto_rawDesc = "" +
	"\n" +
	"\x18proto/filetransfer.proto\x12\ffiletransfer\"\x82\x01\n" +
	"\tFileChunk\x12\x1b\n" +
	"\tfile_name\x18\x01 \x01(\tR\bfileName\x12\x12\n" +
	"\x04data\x18\x02 \x01(\fR\x04data\x12!\n" +
	"\fchunk_number\x18\x03 \x01(\x03R\vchunkNumber\x12!\n" +
	"\ftotal_chunks\x18\x04 \x01(\x03R\vtotalChunks\"q\n" +
	"\x14FileTransferResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x12%\n" +
	"\x0ebytes_received\x18\x03 \x01(\x03R\rbytesReceived2`\n" +
	"\x13FileTransferService\x12I\n" +
	"\bSendFile\x12\x17.filetransfer.FileChunk\x1a\".filetransfer.FileTransferResponse(\x01B\tZ\a./protob\x06proto3"

var (
	file_proto_filetransfer_proto_rawDescOnce sync.Once
	file_proto_filetransfer_proto_rawDescData []byte
)

func file_proto_filetransfer_proto_rawDescGZIP() []byte {
	file_proto_filetransfer_proto_rawDescOnce.Do(func() {
		file_proto_filetransfer_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_filetransfer_proto_rawDesc), len(file_proto_filetransfer_proto_rawDesc)))
	})
	return file_proto_filetransfer_proto_rawDescData
}

var file_proto_filetransfer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_filetransfer_proto_goTypes = []any{
	(*FileChunk)(nil),            // 0: filetransfer.FileChunk
	(*FileTransferResponse)(nil), // 1: filetransfer.FileTransferResponse
}
var file_proto_filetransfer_proto_depIdxs = []int32{
	0, // 0: filetransfer.FileTransferService.SendFile:input_type -> filetransfer.FileChunk
	1, // 1: filetransfer.FileTransferService.SendFile:output_type -> filetransfer.FileTransferResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_filetransfer_proto_init() }
func file_proto_filetransfer_proto_init() {
	if File_proto_filetransfer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_filetransfer_proto_rawDesc), len(file_proto_filetransfer_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_filetransfer_proto_goTypes,
		DependencyIndexes: file_proto_filetransfer_proto_depIdxs,
		MessageInfos:      file_proto_filetransfer_proto_msgTypes,
	}.Build()
	File_proto_filetransfer_proto = out.File
	file_proto_filetransfer_proto_goTypes = nil
	file_proto_filetransfer_proto_depIdxs = nil
}
