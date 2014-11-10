package btrfs

import (
	"fmt"
	"math"
)

func printDirItemType(di *BtrfsDirItem) {
	Type := di.Type

	switch Type {
	case BTRFS_FT_REG_FILE:
		fmt.Printf("FILE")
		break
	case BTRFS_FT_DIR:
		fmt.Printf("DIR")
		break
	case BTRFS_FT_CHRDEV:
		fmt.Printf("CHRDEV")
		break
	case BTRFS_FT_BLKDEV:
		fmt.Printf("BLKDEV")
		break
	case BTRFS_FT_FIFO:
		fmt.Printf("FIFO")
		break
	case BTRFS_FT_SOCK:
		fmt.Printf("SOCK")
		break
	case BTRFS_FT_SYMLINK:
		fmt.Printf("SYMLINK")
		break
	case BTRFS_FT_XATTR:
		fmt.Printf("XATTR")
		break
	default:
		fmt.Printf("%d", Type)
	}
}
func printInodeRefItem() {
	//	struct extent_buffer *eb, struct btrfs_item *item,
	//				struct btrfs_inode_ref *ref)
	//{
	//	u32 total;
	//	u32 cur = 0;
	//	u32 len;
	//	u32 name_len;
	//	u64 index;
	//	char namebuf[BTRFS_NAME_LEN];
	//	total = btrfs_item_size(eb, item);
	//	while(cur < total) {
	//		name_len = btrfs_inode_ref_name_len(eb, ref);
	//		index = btrfs_inode_ref_index(eb, ref);
	//		len = (name_len <= sizeof(namebuf))? name_len: sizeof(namebuf);
	//		read_extent_buffer(eb, namebuf, (unsigned long)(ref + 1), len);
	//		printf("\t\tinode ref index %llu namelen %u name: %.*s\n",
	//		       (unsigned long long)index, name_len, len, namebuf);
	//		len = sizeof(*ref) + name_len;
	//		ref = (struct btrfs_inode_ref *)((char *)ref + len);
	//		cur += len;
	//	}
	//	return 0;
}
func printObjectid(objectid uint64, Type uint8) {
	switch Type {
	case BTRFS_DEV_EXTENT_KEY:
		fmt.Printf("%d", objectid) /* device id */
		return
	case BTRFS_QGROUP_RELATION_KEY:
		fmt.Printf("%d/%d", objectid>>48,
			objectid&((1<<48)-1))
		return
	case BTRFS_UUID_KEY_SUBVOL, BTRFS_UUID_KEY_RECEIVED_SUBVOL:
		fmt.Printf("0x%d", objectid)
		return
	}

	switch objectid {
	case BTRFS_ROOT_TREE_OBJECTID:
		if Type == BTRFS_DEV_ITEM_KEY {
			fmt.Printf("DEV_ITEMS")
		} else {
			fmt.Printf("ROOT_TREE")
		}
		break
	case BTRFS_EXTENT_TREE_OBJECTID:
		fmt.Printf("EXTENT_TREE")
		break
	case BTRFS_CHUNK_TREE_OBJECTID:
		fmt.Printf("CHUNK_TREE")
		break
	case BTRFS_DEV_TREE_OBJECTID:
		fmt.Printf("DEV_TREE")
		break
	case BTRFS_FS_TREE_OBJECTID:
		fmt.Printf("FS_TREE")
		break
	case BTRFS_ROOT_TREE_DIR_OBJECTID:
		fmt.Printf("ROOT_TREE_DIR")
		break
	case BTRFS_CSUM_TREE_OBJECTID:
		fmt.Printf("CSUM_TREE")
		break
	case BTRFS_BALANCE_OBJECTID:
		fmt.Printf("BALANCE")
		break
	case BTRFS_ORPHAN_OBJECTID:
		fmt.Printf("ORPHAN")
		break
	case BTRFS_TREE_LOG_OBJECTID:
		fmt.Printf("TREE_LOG")
		break
	case BTRFS_TREE_LOG_FIXUP_OBJECTID:
		fmt.Printf("LOG_FIXUP")
		break
	case BTRFS_TREE_RELOC_OBJECTID:
		fmt.Printf("TREE_RELOC")
		break
	case BTRFS_DATA_RELOC_TREE_OBJECTID:
		fmt.Printf("DATA_RELOC_TREE")
		break
	case BTRFS_EXTENT_CSUM_OBJECTID:
		fmt.Printf("EXTENT_CSUM")
		break
	case BTRFS_FREE_SPACE_OBJECTID:
		fmt.Printf("FREE_SPACE")
		break
	case BTRFS_FREE_INO_OBJECTID:
		fmt.Printf("FREE_INO")
		break
	case BTRFS_QUOTA_TREE_OBJECTID:
		fmt.Printf("QUOTA_TREE")
		break
	case BTRFS_UUID_TREE_OBJECTID:
		fmt.Printf("UUID_TREE")
		break
	case BTRFS_MULTIPLE_OBJECTIDS:
		fmt.Printf("MULTIPLE")
		break
	case -1 & math.MaxUint64:
		fmt.Printf("-1")
		break
	case BTRFS_FIRST_CHUNK_TREE_OBJECTID:
		if Type == BTRFS_CHUNK_ITEM_KEY {
			fmt.Printf("FIRST_CHUNK_TREE")
			break
		}
		/* fall-thru */
	default:
		fmt.Printf("%d", objectid)
	}
}
func printKeyType(objectid uint64, Type uint8) {
	if Type == 0 && objectid == BTRFS_FREE_SPACE_OBJECTID {
		fmt.Printf("UNTYPED")
		return
	}

	switch Type {
	case BTRFS_INODE_ITEM_KEY:
		fmt.Printf("INODE_ITEM")
		break
	case BTRFS_INODE_REF_KEY:
		fmt.Printf("INODE_REF")
		break
	case BTRFS_INODE_EXTREF_KEY:
		fmt.Printf("INODE_EXTREF")
		break
	case BTRFS_DIR_ITEM_KEY:
		fmt.Printf("DIR_ITEM")
		break
	case BTRFS_DIR_INDEX_KEY:
		fmt.Printf("DIR_INDEX")
		break
	case BTRFS_DIR_LOG_ITEM_KEY:
		fmt.Printf("DIR_LOG_ITEM")
		break
	case BTRFS_DIR_LOG_INDEX_KEY:
		fmt.Printf("DIR_LOG_INDEX")
		break
	case BTRFS_XATTR_ITEM_KEY:
		fmt.Printf("XATTR_ITEM")
		break
	case BTRFS_ORPHAN_ITEM_KEY:
		fmt.Printf("ORPHAN_ITEM")
		break
	case BTRFS_ROOT_ITEM_KEY:
		fmt.Printf("ROOT_ITEM")
		break
	case BTRFS_ROOT_REF_KEY:
		fmt.Printf("ROOT_REF")
		break
	case BTRFS_ROOT_BACKREF_KEY:
		fmt.Printf("ROOT_BACKREF")
		break
	case BTRFS_EXTENT_ITEM_KEY:
		fmt.Printf("EXTENT_ITEM")
		break
	case BTRFS_METADATA_ITEM_KEY:
		fmt.Printf("METADATA_ITEM")
		break
	case BTRFS_TREE_BLOCK_REF_KEY:
		fmt.Printf("TREE_BLOCK_REF")
		break
	case BTRFS_SHARED_BLOCK_REF_KEY:
		fmt.Printf("SHARED_BLOCK_REF")
		break
	case BTRFS_EXTENT_DATA_REF_KEY:
		fmt.Printf("EXTENT_DATA_REF")
		break
	case BTRFS_SHARED_DATA_REF_KEY:
		fmt.Printf("SHARED_DATA_REF")
		break
	case BTRFS_EXTENT_REF_V0_KEY:
		fmt.Printf("EXTENT_REF_V0")
		break
	case BTRFS_CSUM_ITEM_KEY:
		fmt.Printf("CSUM_ITEM")
		break
	case BTRFS_EXTENT_CSUM_KEY:
		fmt.Printf("EXTENT_CSUM")
		break
	case BTRFS_EXTENT_DATA_KEY:
		fmt.Printf("EXTENT_DATA")
		break
	case BTRFS_BLOCK_GROUP_ITEM_KEY:
		fmt.Printf("BLOCK_GROUP_ITEM")
		break
	case BTRFS_CHUNK_ITEM_KEY:
		fmt.Printf("CHUNK_ITEM")
		break
	case BTRFS_DEV_ITEM_KEY:
		fmt.Printf("DEV_ITEM")
		break
	case BTRFS_DEV_EXTENT_KEY:
		fmt.Printf("DEV_EXTENT")
		break
	case BTRFS_BALANCE_ITEM_KEY:
		fmt.Printf("BALANCE_ITEM")
		break
	case BTRFS_DEV_REPLACE_KEY:
		fmt.Printf("DEV_REPLACE_ITEM")
		break
	case BTRFS_STRING_ITEM_KEY:
		fmt.Printf("STRING_ITEM")
		break
	case BTRFS_QGROUP_STATUS_KEY:
		fmt.Printf("BTRFS_STATUS_KEY")
		break
	case BTRFS_QGROUP_RELATION_KEY:
		fmt.Printf("BTRFS_QGROUP_RELATION_KEY")
		break
	case BTRFS_QGROUP_INFO_KEY:
		fmt.Printf("BTRFS_QGROUP_INFO_KEY")
		break
	case BTRFS_QGROUP_LIMIT_KEY:
		fmt.Printf("BTRFS_QGROUP_LIMIT_KEY")
		break
	case BTRFS_DEV_STATS_KEY:
		fmt.Printf("DEV_STATS_ITEM")
		break
	case BTRFS_UUID_KEY_SUBVOL:
		fmt.Printf("BTRFS_UUID_KEY_SUBVOL")
		break
	case BTRFS_UUID_KEY_RECEIVED_SUBVOL:
		fmt.Printf("BTRFS_UUID_KEY_RECEIVED_SUBVOL")
		break
	default:
		fmt.Printf("UNKNOWN.%d", Type)
	}
}
func BtrfsPrintKey(diskKey *BtrfsDiskKey) {
	objectid := diskKey.Objectid
	Type := diskKey.Type
	offset := diskKey.Offset

	fmt.Printf("key (")
	printObjectid(objectid, Type)
	fmt.Printf(" ")
	printKeyType(objectid, Type)
	switch Type {
	case BTRFS_QGROUP_RELATION_KEY, BTRFS_QGROUP_INFO_KEY, BTRFS_QGROUP_LIMIT_KEY:
		fmt.Printf(" %d/%d)", (offset >> 48),
			(offset & ((1 << 48) - 1)))
		break
	case BTRFS_UUID_KEY_SUBVOL, BTRFS_UUID_KEY_RECEIVED_SUBVOL:
		fmt.Printf(" 0x%d)", offset)
		break
	default:
		if offset == -1&math.MaxUint64 {
			fmt.Printf(" -1)")
		} else {
			fmt.Printf(" %d)", offset)
		}
		break
	}
}

//func btrfsPrintLeaf(root  *BtrfsRoot, l *ExtentBuffer ) {
////	int i
////	char *str
////	struct btrfs_item *item
////	struct btrfs_dir_item *di
////	struct btrfs_inode_item *ii
////	struct btrfs_file_extent_item *fi
////	struct btrfs_block_group_item *bi
////	struct btrfs_extent_data_ref *dref
////	struct btrfs_shared_data_ref *sref
////	struct btrfs_inode_ref *iref
////	struct btrfs_inode_extref *iref2
////	struct btrfs_dev_extent *dev_extent
////	struct btrfs_disk_key disk_key
////	struct btrfs_block_group_item bg_item
////	struct btrfs_dir_log_item *dlog
////	struct btrfs_qgroup_info_item *qg_info
////	struct btrfs_qgroup_limit_item *qg_limit
////	struct btrfs_qgroup_status_item *qg_status
////	u32 nr = btrfs_header_nritems(l)
////	u64 objectid
////	u32 type
////	char bg_flags_str[32]
//
////	fmt.Printf("leaf %d items %d free space %d generation %d owner %d\n",
////		btrfs_header_bytenr(l), nr,
////		btrfs_leaf_free_space(root, l),
////		btrfs_header_generation(l),
////		btrfs_header_owner(l))
////	print_uuids(l)
////	fflush(stdout)
//	nr:=0
//	for i := 0 ; i < nr;  i++ {
////		item = btrfs_item_nr(i)
////		btrfs_item_key(l, &disk_key, i)
////		objectid = btrfs_disk_key_objectid(&disk_key)
////		Type = btrfs_disk_key_type(&disk_key)
//		fmt.Printf("\titem %d ", i)
//		btrfs_print_key(&disk_key)
//		fmt.Printf(" itemoff %d itemsize %d\n",
//			btrfs_item_offset(l, item),
//			btrfs_item_size(l, item))
//
//		if (Type == 0 && objectid == BTRFS_FREE_SPACE_OBJECTID) {
//			print_free_space_header(l, i)
//			}
//
//		switch (Type) {
//		case BTRFS_INODE_ITEM_KEY:
////			ii = btrfs_item_ptr(l, i, struct btrfs_inode_item)
//			fmt.Printf("\t\tinode generation %d transid %d size %d block group %d mode %o links %u uid %u gid %u rdev %d flags 0x%llx\n",
//			       btrfs_inode_generation(l, ii),
//			       btrfs_inode_transid(l, ii),
//			       btrfs_inode_size(l, ii),
//			       btrfs_inode_block_group(l,ii),
//			       btrfs_inode_mode(l, ii),
//			       btrfs_inode_nlink(l, ii),
//			       btrfs_inode_uid(l, ii),
//			       btrfs_inode_gid(l, ii),
//			       btrfs_inode_rdev(l,ii),
//			       btrfs_inode_flags(l,ii))
//			break
//		case BTRFS_INODE_REF_KEY:
////			iref = btrfs_item_ptr(l, i, struct btrfs_inode_ref)
//			print_inode_ref_item(l, item, iref)
//			break
//		case BTRFS_INODE_EXTREF_KEY:
////			iref2 = btrfs_item_ptr(l, i, struct btrfs_inode_extref)
//			print_inode_extref_item(l, item, iref2)
//			break
//		case BTRFS_DIR_ITEM_KEY:
//		case BTRFS_DIR_INDEX_KEY:
//		case BTRFS_XATTR_ITEM_KEY:
////			di = btrfs_item_ptr(l, i, struct btrfs_dir_item)
//			print_dir_item(l, item, di)
//			break
//		case BTRFS_DIR_LOG_INDEX_KEY:
//		case BTRFS_DIR_LOG_ITEM_KEY:
////			dlog = btrfs_item_ptr(l, i, struct btrfs_dir_log_item)
//			fmt.Printf("\t\tdir log end %Lu\n",
//			       btrfs_dir_log_end(l, dlog))
//		       break
//		case BTRFS_ORPHAN_ITEM_KEY:
//			fmt.Printf("\t\torphan item\n")
//			break
//		case BTRFS_ROOT_ITEM_KEY:
//			print_root(l, i)
//			break
//		case BTRFS_ROOT_REF_KEY:
//			print_root_ref(l, i, "ref")
//			break
//		case BTRFS_ROOT_BACKREF_KEY:
//			print_root_ref(l, i, "backref")
//			break
//		case BTRFS_EXTENT_ITEM_KEY:
//			print_extent_item(l, i, 0)
//			break
//		case BTRFS_METADATA_ITEM_KEY:
//			print_extent_item(l, i, 1)
//			break
//		case BTRFS_TREE_BLOCK_REF_KEY:
//			fmt.Printf("\t\ttree block backref\n")
//			break
//		case BTRFS_SHARED_BLOCK_REF_KEY:
//			fmt.Printf("\t\tshared block backref\n")
//			break
//		case BTRFS_EXTENT_DATA_REF_KEY:
////			dref = btrfs_item_ptr(l, i, struct btrfs_extent_data_ref)
//			fmt.Printf("\t\textent data backref root %d objectid %d offset %d count %u\n",
//			       btrfs_extent_data_ref_root(l, dref),
//			       btrfs_extent_data_ref_objectid(l, dref),
//			       btrfs_extent_data_ref_offset(l, dref),
//			       btrfs_extent_data_ref_count(l, dref))
//			break
//		case BTRFS_SHARED_DATA_REF_KEY:
////			sref = btrfs_item_ptr(l, i, struct btrfs_shared_data_ref)
//			fmt.Printf("\t\tshared data backref count %u\n",
//			       btrfs_shared_data_ref_count(l, sref))
//			break
//		case BTRFS_EXTENT_REF_V0_KEY:
//
//			print_extent_ref_v0(l, i)
//
//			break
//		case BTRFS_CSUM_ITEM_KEY:
//			fmt.Printf("\t\tcsum item\n")
//			break
//		case BTRFS_EXTENT_CSUM_KEY:
//			fmt.Printf("\t\textent csum item\n")
//			break
//		case BTRFS_EXTENT_DATA_KEY:
////			fi = btrfs_item_ptr(l, i,
////					    struct btrfs_file_extent_item)
//			print_file_extent_item(l, item, i, fi)
//			break
//		case BTRFS_BLOCK_GROUP_ITEM_KEY:
////			bi = btrfs_item_ptr(l, i,
////					    struct btrfs_block_group_item)
////			read_extent_buffer(l, &bg_item, (unsigned long)bi,
////					   sizeof(bg_item))
//			memset(bg_flags_str, 0, sizeof(bg_flags_str))
//			bg_flags_to_str(btrfs_block_group_flags(&bg_item),
//					bg_flags_str)
//			fmt.Printf("\t\tblock group used %d chunk_objectid %d flags %s\n",
//			       btrfs_block_group_used(&bg_item),
//			       btrfs_block_group_chunk_objectid(&bg_item),
//			       bg_flags_str)
//			break
//		case BTRFS_CHUNK_ITEM_KEY:
////			print_chunk(l, btrfs_item_ptr(l, i, struct btrfs_chunk))
//			break
//		case BTRFS_DEV_ITEM_KEY:
////			print_dev_item(l, btrfs_item_ptr(l, i,
////					struct btrfs_dev_item))
//			break
//		case BTRFS_DEV_EXTENT_KEY:
////			dev_extent = btrfs_item_ptr(l, i,
////						    struct btrfs_dev_extent)
//			fmt.Printf("\t\tdev extent chunk_tree %d\n\t\tchunk objectid %d chunk offset %d length %d\n",
//
//			       btrfs_dev_extent_chunk_tree(l, dev_extent),
//
//			       btrfs_dev_extent_chunk_objectid(l, dev_extent),
//
//			       btrfs_dev_extent_chunk_offset(l, dev_extent),
//
//			       btrfs_dev_extent_length(l, dev_extent))
//			break
//		case BTRFS_QGROUP_STATUS_KEY:
////			qg_status = btrfs_item_ptr(l, i,
////					struct btrfs_qgroup_status_item)
//			fmt.Printf("\t\tversion %d generation %d flags %#llx scan %lld\n",
//
//				btrfs_qgroup_status_version(l, qg_status),
//
//				btrfs_qgroup_status_generation(l, qg_status),
//
//				btrfs_qgroup_status_flags(l, qg_status),
//
//				btrfs_qgroup_status_scan(l, qg_status))
//			break
//		case BTRFS_QGROUP_RELATION_KEY:
//			break
//		case BTRFS_QGROUP_INFO_KEY:
////			qg_info = btrfs_item_ptr(l, i,
////						 struct btrfs_qgroup_info_item)
//			fmt.Printf("\t\tgeneration %d\n\t\treferenced %d referenced compressed %d\n\t\texclusive %d exclusive compressed %d\n",
//
//			       btrfs_qgroup_info_generation(l, qg_info),
//
//			       btrfs_qgroup_info_referenced(l, qg_info),
//
//			       btrfs_qgroup_info_referenced_compressed(l,
//								       qg_info),
//
//			       btrfs_qgroup_info_exclusive(l, qg_info),
//
//			       btrfs_qgroup_info_exclusive_compressed(l,
//								      qg_info))
//			break
//		case BTRFS_QGROUP_LIMIT_KEY:
////			qg_limit = btrfs_item_ptr(l, i,
////					 struct btrfs_qgroup_limit_item)
//			fmt.Printf("\t\tflags %llx\n\t\tmax referenced %lld max exclusive %lld\n\t\trsv referenced %lld rsv exclusive %lld\n",
//
//			       btrfs_qgroup_limit_flags(l, qg_limit),
//
//			       btrfs_qgroup_limit_max_referenced(l, qg_limit),
//
//			       btrfs_qgroup_limit_max_exclusive(l, qg_limit),
//
//			       btrfs_qgroup_limit_rsv_referenced(l, qg_limit),
//
//			       btrfs_qgroup_limit_rsv_exclusive(l, qg_limit))
//			break
//		case BTRFS_UUID_KEY_SUBVOL:
//		case BTRFS_UUID_KEY_RECEIVED_SUBVOL:
//			print_uuid_item(l, btrfs_item_ptr_offset(l, i),
//					btrfs_item_size_nr(l, i))
//			break
//		case BTRFS_STRING_ITEM_KEY:
//			/* dirty, but it's simple */
////			str = l->data + btrfs_item_ptr_offset(l, i)
//			fmt.Printf("\t\titem data %.*s\n", btrfs_item_size(l, item), str)
//			break
//		case BTRFS_DEV_STATS_KEY:
//			fmt.Printf("\t\tdevice stats\n")
//			break
//		}
//		fflush(stdout)
//	}
//}
