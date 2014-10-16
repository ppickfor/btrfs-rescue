package btrfs

import ()

type rb_root struct {
}
type cache_tree struct {
	root rb_root
}

type list_head struct {
	next *list_head
	prev *list_head
}

type block_group_tree struct {
	tree         cache_tree
	block_groups list_head
}
type device_extent_tree struct {
	tree cache_tree
	/*
	 * The idea is:
	 * When checking the chunk information, we move the device extents
	 * that has its chunk to the chunk's device extents list. After the
	 * check, if there are still some device extents in no_chunk_orphans,
	 * it means there are some device extents which don't belong to any
	 * chunk.
	 *
	 * The usage of no_device_orphans is the same as the first one, but it
	 * is for the device information check.
	 */
	no_chunk_orphans  list_head
	no_device_orphans list_head
}
