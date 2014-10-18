// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs=true -debug-define=true btrfs_h.go

package btrfs

type Mutex struct {
	Lock uint64
}
type Rb_node struct {
	Next *Rb_node
}
type Rb_root struct {
	Next *Rb_node
}
type Cache_extent struct {
	Node		Rb_node
	Objectid	uint64
	Start		uint64
	Size		uint64
}
type Cache_tree struct {
	Root Rb_root
}
type List_head struct {
	Next	*List_head
	Prev	*List_head
}
type Extent_io_tree struct {
	State	Cache_tree
	Cache	Cache_tree
	Lru	List_head
	Size	uint64
}
type Block_group_record struct {
	Cache		Cache_extent
	List		List_head
	Generation	uint64
	Objectid	uint64
	Type		uint8
	Pad_cgo_0	[7]byte
	Offset		uint64
	Flags		uint64
}
type Block_group_tree struct {
	Tree	Cache_tree
	Groups	List_head
}
type Device_record struct {
	Node		Rb_node
	Devid		uint64
	Generation	uint64
	Objectid	uint64
	Type		uint8
	Pad_cgo_0	[7]byte
	Offset		uint64
	Total_byte	uint64
	Byte_used	uint64
	Real_used	uint64
}
type Stripe struct {
	Devid	uint64
	Offset	uint64
	Uuid	[16]uint8
}
type Chunk_record struct {
	Cache		Cache_extent
	List		List_head
	Dextents	List_head
	Bg_rec		*Block_group_record
	Generation	uint64
	Objectid	uint64
	Type		uint8
	Pad_cgo_0	[7]byte
	Offset		uint64
	Owner		uint64
	Length		uint64
	Type_flags	uint64
	Stripe_len	uint64
	Num_stripes	uint16
	Sub_stripes	uint16
	Pad_cgo_1	[4]byte
	Io_align	uint64
	Io_width	uint64
	Sector_size	uint64
	Stripes		[0]Stripe
}
type Device_extent_record struct {
	Cache		Cache_extent
	Chunk_list	List_head
	Device_list	List_head
	Generation	uint64
	Objectid	uint64
	Type		uint8
	Pad_cgo_0	[7]byte
	Offset		uint64
	Chunk_objecteid	uint64
	Chunk_offset	uint64
	Length		uint64
}
type Device_extent_tree struct {
	Tree		Cache_tree
	Chunk_orphans	List_head
	Device_orphans	List_head
}
type Root_info struct {
	Rb_node		Rb_node
	Sort_node	Rb_node
	Root_id		uint64
	Root_offset	uint64
	Flags		uint64
	Ref_tree	uint64
	Dir_id		uint64
	Top_id		uint64
	Gen		uint64
	Ogen		uint64
	Otime		uint64
	Uuid		[16]uint8
	Puuid		[16]uint8
	Ruuid		[16]uint8
	Path		*int8
	Name		*int8
	Full_path	*int8
	Deleted		int32
	Pad_cgo_0	[4]byte
}
type Cmd_struct struct {
	Token		*int8
	Fn		*[0]byte
	Usagestr	**int8
	Next		*Cmd_group
	Hidden		int32
	Pad_cgo_0	[4]byte
}
type Cmd_group struct {
	Usagestr	**int8
	Infostr		*int8
	Commands	[0]byte
}
type Btrfs_disk_key struct {
	Objectid	uint64
	Type		uint8
	Pad_cgo_0	[7]byte
	Offset		uint64
}
type Btrfs_key struct {
	Objectid	uint64
	Type		uint8
	Pad_cgo_0	[7]byte
	Offset		uint64
}
type Btrfs_mapping_tree struct {
	Tree Cache_tree
}
type Btrfs_dev_item struct {
	Devid		uint64
	Total_bytes	uint64
	Bytes_used	uint64
	Io_align	uint64
	Io_width	uint64
	Sector_size	uint64
	Type		uint64
	Generation	uint64
	Start_offset	uint64
	Dev_group	uint64
	Seek_speed	uint8
	Bandwidth	uint8
	Uuid		[16]uint8
	Fsid		[16]uint8
	Pad_cgo_0	[6]byte
}
type Btrfs_stripe struct {
	Devid	uint64
	Offset	uint64
	Uuid	[16]uint8
}
type Btrfs_chunk struct {
	Length		uint64
	Owner		uint64
	Stripe_len	uint64
	Type		uint64
	Io_align	uint64
	Io_width	uint64
	Sector_size	uint64
	Num_stripes	uint16
	Sub_stripes	uint16
	Pad_cgo_0	[4]byte
	Stripe		Btrfs_stripe
}
type Btrfs_free_space_entry struct {
	Offset		uint64
	Bytes		uint64
	Type		uint8
	Pad_cgo_0	[7]byte
}
type Btrfs_free_space_header struct {
	Location	Btrfs_disk_key
	Generation	uint64
	Entries		uint64
	Bitmaps		uint64
}
type Btrfs_header struct {
	Csum		[32]uint8
	Fsid		[16]uint8
	Bytenr		uint64
	Flags		uint64
	Tree_uuid	[16]uint8
	Generation	uint64
	Owner		uint64
	Nritems		uint64
	Level		uint8
	Pad_cgo_0	[7]byte
}
type Btrfs_root_backup struct {
	Tree_root		uint64
	Tree_root_gen		uint64
	Chunk_root		uint64
	Chunk_root_gen		uint64
	Extent_root		uint64
	Extent_root_gen		uint64
	Fs_root			uint64
	Fs_root_gen		uint64
	Dev_root		uint64
	Dev_root_gen		uint64
	Csum_root		uint64
	Csum_root_gen		uint64
	Total_bytes		uint64
	Bytes_used		uint64
	Num_devices		uint64
	Unsed_64		[4]uint64
	Tree_root_level		uint8
	Chunk_root_level	uint8
	Extent_root_level	uint8
	Fs_root_level		uint8
	Dev_root_level		uint8
	Csum_root_level		uint8
	Unused_8		[10]uint8
}
type Btrfs_super_block struct {
	Csum			[32]uint8
	Fsid			[16]uint8
	Bytenr			uint64
	Flags			uint64
	Magic			uint64
	Generation		uint64
	Root			uint64
	Chunk_root		uint64
	Log_root		uint64
	Log_root_transid	uint64
	Total_bytes		uint64
	Bytes_used		uint64
	Root_dir_objectid	uint64
	Num_devices		uint64
	Sectorsize		uint64
	Nodesize		uint64
	Leafsize		uint64
	Stripesize		uint64
	Sys_chunk_array_size	uint64
	Chunk_root_generation	uint64
	Compat_flags		uint64
	Compat_ro_flags		uint64
	Incompat_flags		uint64
	Csum_type		uint16
	Root_level		uint8
	Chunk_root_level	uint8
	Log_root_level		uint8
	Pad_cgo_0		[3]byte
	Dev_item		Btrfs_dev_item
	Label			[256]int8
	Cache_generation	uint64
	Uuid_tree_generation	uint64
	Reserved		[30]uint64
	Sys_chunk_array		[2048]uint8
	Super_roots		[4]Btrfs_root_backup
}
type Btrfs_item struct {
	Key	Btrfs_disk_key
	Offset	uint64
	Size	uint64
}
type Btrfs_leaf struct {
	Header	Btrfs_header
	Items	[0]Btrfs_item
}
type Btrfs_key_ptr struct {
	Key		Btrfs_disk_key
	Blockptr	uint64
	Generation	uint64
}
type Btrfs_node struct {
	Header	Btrfs_header
	Ptrs	[0]Btrfs_key_ptr
}
type Btrfs_path struct {
	Nodes		[8]*Extent_buffer
	Slots		[8]int32
	Locks		[8]int32
	Reada		int32
	Level		int32
	Pad_cgo_0	[8]byte
}
type Btrfs_extent_item struct {
	Refs		uint64
	Generation	uint64
	Flags		uint64
}
type Btrfs_extent_item_v0 struct {
	Refs uint64
}
type Btrfs_tree_block_info struct {
	Key		Btrfs_disk_key
	Level		uint8
	Pad_cgo_0	[7]byte
}
type Btrfs_extent_data_ref struct {
	Root		uint64
	Objectid	uint64
	Offset		uint64
	Count		uint64
}
type Btrfs_shared_data_ref struct {
	Count uint64
}
type Btrfs_extent_inline_ref struct {
	Type		uint8
	Pad_cgo_0	[7]byte
	Offset		uint64
}
type Btrfs_extent_ref_v0 struct {
	Root		uint64
	Generation	uint64
	Objectid	uint64
	Count		uint64
}
type Btrfs_dev_extent struct {
	Tree		uint64
	Objectid	uint64
	Offset		uint64
	Length		uint64
	Tree_uuid	[16]uint8
}
type Btrfs_inode_ref struct {
	Index		uint64
	Len		uint16
	Pad_cgo_0	[6]byte
}
type Btrfs_inode_extref struct {
	Parent_objectid	uint64
	Index		uint64
	Name_len	uint16
	Name		[0]byte
	Pad_cgo_0	[6]byte
}
type Btrfs_timespec struct {
	Sec	uint64
	Nsec	uint64
}
type Btrfs_inode_item struct {
	Generation	uint64
	Transid		uint64
	Size		uint64
	Nbytes		uint64
	Group		uint64
	Nlink		uint64
	Uid		uint64
	Gid		uint64
	Mode		uint64
	Rdev		uint64
	Flags		uint64
	Sequence	uint64
	Reserved	[4]uint64
	Atime		Btrfs_timespec
	Ctime		Btrfs_timespec
	Mtime		Btrfs_timespec
	Otime		Btrfs_timespec
}
type Btrfs_dir_log_item struct {
	End uint64
}
type Btrfs_dir_item struct {
	Location	Btrfs_disk_key
	Transid		uint64
	Data_len	uint16
	Name_len	uint16
	Type		uint8
	Pad_cgo_0	[3]byte
}
type Btrfs_root_item_v0 struct {
	Inode		Btrfs_inode_item
	Generation	uint64
	Root_dirid	uint64
	Bytenr		uint64
	Byte_limit	uint64
	Bytes_used	uint64
	Last_snapshot	uint64
	Flags		uint64
	Refs		uint64
	Drop_progress	Btrfs_disk_key
	Drop_level	uint8
	Level		uint8
	Pad_cgo_0	[6]byte
}
type Btrfs_root_item struct {
	Inode		Btrfs_inode_item
	Generation	uint64
	Root_dirid	uint64
	Bytenr		uint64
	Byte_limit	uint64
	Bytes_used	uint64
	Last_snapshot	uint64
	Flags		uint64
	Refs		uint64
	Drop_progress	Btrfs_disk_key
	Drop_level	uint8
	Level		uint8
	Pad_cgo_0	[6]byte
	Generation_v2	uint64
	Uuid		[16]uint8
	Parent_uuid	[16]uint8
	Received_uuid	[16]uint8
	Ctransid	uint64
	Otransid	uint64
	Stransid	uint64
	Rtransid	uint64
	Ctime		Btrfs_timespec
	Otime		Btrfs_timespec
	Stime		Btrfs_timespec
	Rtime		Btrfs_timespec
	Reserved	[8]uint64
}
type Btrfs_root_ref struct {
	Dirid		uint64
	Sequence	uint64
	Len		uint16
	Pad_cgo_0	[6]byte
}
type Btrfs_file_extent_item struct {
	Generation	uint64
	Ram_bytes	uint64
	Compression	uint8
	Encryption	uint8
	Other_encoding	uint16
	Type		uint8
	Pad_cgo_0	[3]byte
	Disk_bytenr	uint64
	Disk_num_bytes	uint64
	Offset		uint64
	Num_bytes	uint64
}
type Btrfs_csum_item struct {
	Csum uint8
}
type Btrfs_qgroup_status_item struct {
	Version		uint64
	Generation	uint64
	Flags		uint64
	Scan		uint64
}
type Btrfs_block_group_item struct {
	Used		uint64
	Objectid	uint64
	Flags		uint64
}
type Btrfs_qgroup_info_item struct {
	Generation		uint64
	Referenced		uint64
	Referenced_compressed	uint64
	Exclusive		uint64
	Exclusive_compressed	uint64
}
type Btrfs_qgroup_limit_item struct {
	Flags		uint64
	Max_referenced	uint64
	Max_exclusive	uint64
	Rsv_referenced	uint64
	Rsv_exclusive	uint64
}
type Btrfs_space_info struct {
	Flags		uint64
	Total_bytes	uint64
	Bytes_used	uint64
	Bytes_pinned	uint64
	Full		int32
	Pad_cgo_0	[4]byte
	List		List_head
}
type Btrfs_block_group_cache struct {
	Cache		Cache_extent
	Key		Btrfs_key
	Item		Btrfs_block_group_item
	Space_info	*Btrfs_space_info
	Free_space_ctl	*Btrfs_free_space_ctl
	Pinned		uint64
	Flags		uint64
	Cached		int32
	Ro		int32
}
type Btrfs_fs_info struct {
	Fsid				[16]uint8
	Chunk_tree_uuid			[16]uint8
	Fs_root				*Btrfs_root
	Extent_root			*Btrfs_root
	Tree_root			*Btrfs_root
	Chunk_root			*Btrfs_root
	Dev_root			*Btrfs_root
	Csum_root			*Btrfs_root
	Quota_root			*Btrfs_root
	Fs_root_tree			Rb_root
	Log_root_tree			*Btrfs_root
	Extent_cache			Extent_io_tree
	Free_space_cache		Extent_io_tree
	Block_group_cache		Extent_io_tree
	Pinned_extents			Extent_io_tree
	Pending_del			Extent_io_tree
	Extent_ins			Extent_io_tree
	Mapping_tree			Btrfs_mapping_tree
	Generation			uint64
	Last_trans_committed		uint64
	Avail_data_alloc_bits		uint64
	Avail_metadata_alloc_bits	uint64
	Avail_system_alloc_bits		uint64
	Data_alloc_profile		uint64
	Metadata_alloc_profile		uint64
	System_alloc_profile		uint64
	Alloc_start			uint64
	Running_transaction		*Btrfs_trans_handle
	Super_copy			*Btrfs_super_block
	Fs_mutex			Mutex
	Super_bytenr			uint64
	Total_pinned			uint64
	Extent_ops			*Btrfs_extent_ops
	Dirty_cowonly_roots		List_head
	Recow_ebs			List_head
	Fs_devices			*Btrfs_fs_devices
	Space_info			List_head
	System_allocs			int32
	Pad_cgo_0			[4]byte
	Free_extent_hook		*[0]byte
	Fsck_extent_cache		*Cache_tree
	Corrupt_blocks			*Cache_tree
}
type Btrfs_root struct {
	Node			*Extent_buffer
	Commit_root		*Extent_buffer
	Root_item		Btrfs_root_item
	Root_key		Btrfs_key
	Fs_info			*Btrfs_fs_info
	Objectid		uint64
	Last_trans		uint64
	Sectorsize		uint64
	Nodesize		uint64
	Leafsize		uint64
	Stripesize		uint64
	Ref_cows		int32
	Track_dirty		int32
	Type			uint64
	Highest_inode		uint64
	Last_inode_alloc	uint64
	Dirty_list		List_head
	Rb_node			Rb_node
}
type Extent_state struct {
	Node		Cache_extent
	Start		uint64
	End		uint64
	Refs		int32
	Pad_cgo_0	[4]byte
	State		uint64
	Xprivate	uint64
}
type Extent_buffer struct {
	Cache_node	Cache_extent
	Start		uint64
	Dev_bytenr	uint64
	Len		uint64
	Tree		*Extent_io_tree
	Lru		List_head
	Recow		List_head
	Refs		int32
	Flags		int32
	Fd		int32
	Data		[0]byte
	Pad_cgo_0	[4]byte
}
type Btrfs_free_space struct {
	Index	Rb_node
	Offset	uint64
	Bytes	uint64
	Bitmap	*uint64
	List	List_head
}
type Btrfs_free_space_ctl struct {
	Free_space_offset	Rb_root
	Free_space		uint64
	Extents_thresh		int32
	Free_extents		int32
	Total_bitmaps		int32
	Unit			int32
	Start			uint64
	Private			*byte
}
type Btrfs_ioctl_vol_args struct {
	Fd	int64
	Name	[4088]int8
}
type Btrfs_qgroup_limit struct {
	Flags		uint64
	Max_referenced	uint64
	Max_exclusive	uint64
	Rsv_referenced	uint64
	Rsv_exclusive	uint64
}
type Btrfs_qgroup_inherit struct {
	Flags		uint64
	Num_groups		uint64
	Num_ref_copies	uint64
	Num_excl_copies	uint64
	Lim		Btrfs_qgroup_limit
	Qgroups		[0]uint64
}
type Btrfs_ioctl_qgroup_limit_args struct {
	Qgroupid	uint64
	Lim		Btrfs_qgroup_limit
}
type Btrfs_ioctl_vol_args_v2 struct {
	Fd	int64
	Transid	uint64
	Flags	uint64
	Anon0	[32]byte
	Name	[4040]int8
}
type Btrfs_scrub_progress struct {
	Data_extents_scrubbed	uint64
	Tree_extents_scrubbed	uint64
	Data_bytes_scrubbed	uint64
	Tree_bytes_scrubbed	uint64
	Read_errors		uint64
	Csum_errors		uint64
	Verify_errors		uint64
	No_csum			uint64
	Csum_discards		uint64
	Super_errors		uint64
	Malloc_errors		uint64
	Uncorrectable_errors	uint64
	Corrected_errors	uint64
	Last_physical		uint64
	Unverified_errors	uint64
}
type Btrfs_ioctl_scrub_args struct {
	Devid		uint64
	Start		uint64
	End		uint64
	Flags		uint64
	Progress	Btrfs_scrub_progress
	Unused		[109]uint64
}
type Btrfs_ioctl_dev_replace_start_params struct {
	Srcdevid			uint64
	Cont_reading_from_srcdev_mode	uint64
	Srcdev_name			[1025]uint8
	Tgtdev_name			[1025]uint8
	Pad_cgo_0			[6]byte
}
type Btrfs_ioctl_dev_replace_status_params struct {
	Replace_state			uint64
	Progress_1000			uint64
	Time_started			uint64
	Time_stopped			uint64
	Num_write_errors		uint64
	Num_uncorrectable_read_errors	uint64
}
type Btrfs_ioctl_dev_replace_args struct {
	Cmd	uint64
	Result	uint64
	Anon0	[2072]byte
	Spare	[64]uint64
}
type Btrfs_ioctl_dev_info_args struct {
	Devid		uint64
	Uuid		[16]uint8
	Bytes_used	uint64
	Total_bytes	uint64
	Unused		[379]uint64
	Path		[1024]uint8
}
type Btrfs_ioctl_fs_info_args struct {
	Max_id		uint64
	Num_devices	uint64
	Fsid		[16]uint8
	Reserved	[124]uint64
}
type Btrfs_balance_args struct {
	Profiles	uint64
	Usage		uint64
	Devid		uint64
	Pstart		uint64
	Pend		uint64
	Vstart		uint64
	Vend		uint64
	Target		uint64
	Flags		uint64
	Limit		uint64
	Unused		[7]uint64
}
type Btrfs_balance_progress struct {
	Expected	uint64
	Considered	uint64
	Completed	uint64
}
type Btrfs_ioctl_balance_args struct {
	Flags	uint64
	State	uint64
	Data	Btrfs_balance_args
	Meta	Btrfs_balance_args
	Sys	Btrfs_balance_args
	Stat	Btrfs_balance_progress
	Unused	[72]uint64
}
type Btrfs_ioctl_search_key struct {
	Tree_id		uint64
	Min_objectid	uint64
	Max_objectid	uint64
	Min_offset	uint64
	Max_offset	uint64
	Min_transid	uint64
	Max_transid	uint64
	Min_type	uint64
	Max_type	uint64
	Nr_items	uint64
	Unused		uint64
	Unused1		uint64
	Unused2		uint64
	Unused3		uint64
	Unused4		uint64
}
type Btrfs_ioctl_search_header struct {
	Transid		uint64
	Objectid	uint64
	Offset		uint64
	Type		uint64
	Len		uint64
}
type Btrfs_ioctl_search_args struct {
	Key	Btrfs_ioctl_search_key
	Buf	[3976]int8
}
type Btrfs_ioctl_ino_lookup_args struct {
	Treeid		uint64
	Objectid	uint64
	Name		[4080]int8
}
type Btrfs_ioctl_defrag_range_args struct {
	Start		uint64
	Len		uint64
	Flags		uint64
	Extent_thresh	uint64
	Compress_type	uint64
	Unused		[4]uint64
}
type Btrfs_ioctl_space_info struct {
	Flags		uint64
	Total_bytes	uint64
	Used_bytes	uint64
}
type Btrfs_ioctl_space_args struct {
	Space_slots	uint64
	Total_spaces	uint64
	Spaces		[0]Btrfs_ioctl_space_info
}
type Btrfs_data_container struct {
	Bytes_left	uint64
	Bytes_missing	uint64
	Elem_cnt	uint64
	Elem_missed	uint64
	Val		[0]uint64
}
type Btrfs_ioctl_ino_path_args struct {
	Inum		uint64
	Size		uint64
	Reserved	[4]uint64
	Fspath		uint64
}
type Btrfs_ioctl_logical_ino_args struct {
	Logical		uint64
	Size		uint64
	Reserved	[4]uint64
	Inodes		uint64
}
type Btrfs_ioctl_timespec struct {
	Sec	uint64
	Nsec	uint64
}
type Btrfs_ioctl_received_subvol_args struct {
	Uuid		[16]int8
	Stransid	uint64
	Rtransid	uint64
	Stime		Btrfs_ioctl_timespec
	Rtime		Btrfs_ioctl_timespec
	Flags		uint64
	Reserved	[16]uint64
}
type Btrfs_ioctl_send_args struct {
	Send_fd			int64
	Clone_sources_count	uint64
	Clone_sources		*uint64
	Parent_root		uint64
	Flags			uint64
	Reserved		[4]uint64
}
type Btrfs_ioctl_get_dev_stats struct {
	Devid	uint64
	Items	uint64
	Flags	uint64
	Values	[5]uint64
	Unused	[121]uint64
}
type Btrfs_ioctl_quota_ctl_args struct {
	Cmd	uint64
	Status	uint64
}
type Btrfs_ioctl_quota_rescan_args struct {
	Flags		uint64
	Progress	uint64
	Reserved	[6]uint64
}
type Btrfs_ioctl_qgroup_assign_args struct {
	Assign	uint64
	Src	uint64
	Dst	uint64
}
type Btrfs_ioctl_qgroup_create_args struct {
	Create		uint64
	Qgroupid	uint64
}
type Btrfs_ioctl_clone_range_args struct {
	Src_fd		int64
	Src_offset	uint64
	Src_length	uint64
	Dest_offset	uint64
}
type Vma_shared struct {
	Tree_node int32
}
type Vm_area_struct struct {
	Pgoff		uint64
	Start		uint64
	End		uint64
	Shared		Vma_shared
	Pad_cgo_0	[4]byte
}
type Page struct {
	Index uint64
}
type Radix_tree_root struct{}
type Btrfs_corrupt_block struct {
	Cache		Cache_extent
	Key		Btrfs_key
	Level		int32
	Pad_cgo_0	[4]byte
}
type Btrfs_stream_header struct {
	Magic		[13]int8
	Pad_cgo_0	[3]byte
	Version		uint64
}
type Btrfs_cmd_header struct {
	Len		uint64
	Cmd		uint16
	Pad_cgo_0	[6]byte
	Crc		uint64
}
type Btrfs_tlv_header struct {
	Type	uint16
	Len	uint16
}
type Subvol_info struct {
	Root_id		uint64
	Uuid		[16]uint8
	Parent_uuid	[16]uint8
	Received_uuid	[16]uint8
	Ctransid	uint64
	Otransid	uint64
	Stransid	uint64
	Rtransid	uint64
	Path		*int8
}
type Subvol_uuid_search struct {
	Fd int32
}
type Btrfs_trans_handle struct {
	Transid			uint64
	Alloc_exclude_start	uint64
	Alloc_exclude_nr	uint64
	Blocks_reserved		uint64
	Blocks_used		uint64
	Block_group		*Btrfs_block_group_cache
}
type Ulist_iterator struct {
	List *List_head
}
type Ulist_node struct {
	Val	uint64
	Aux	uint64
	List	List_head
	Node	Rb_node
}
type Ulist struct {
	Nnodes	uint64
	Nodes	List_head
	Root	Rb_root
}
type Btrfs_device struct {
	Dev_list		List_head
	Dev_root		*Btrfs_root
	Fs_devices		*Btrfs_fs_devices
	Total_ios		uint64
	Fd			int32
	Writeable		int32
	Name			*int8
	Label			*int8
	Total_devs		uint64
	Super_bytes_used	uint64
	Devid			uint64
	Total_bytes		uint64
	Bytes_used		uint64
	Io_align		uint64
	Io_width		uint64
	Sector_size		uint64
	Type			uint64
	Uuid			[16]uint8
}
type Btrfs_fs_devices struct {
	Fsid		[16]uint8
	Latest_devid	uint64
	Latest_trans	uint64
	Lowest_devid	uint64
	Latest_bdev	int32
	Lowest_bdev	int32
	Devices		List_head
	List		List_head
	Seeding		int32
	Pad_cgo_0	[4]byte
	Seed		*Btrfs_fs_devices
}
type Btrfs_bio_stripe struct {
	Dev		*Btrfs_device
	Physical	uint64
}
type Btrfs_multi_bio struct {
	Error	int32
	Num_stripes	int32
	Stripes	[0]Btrfs_bio_stripe
}
type Map_lookup struct {
	Ce		Cache_extent
	Type		uint64
	Io_align	int32
	Io_width	int32
	Stripe_len	int32
	Sector_size	int32
	Num_stripes	int32
	Sub_stripes	int32
	Stripes		[0]Btrfs_bio_stripe
}
type Btrfs_extent_ops struct {
	Alloc_extent	*[0]byte
	Free_extent	*[0]byte
}
