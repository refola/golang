// raidn.go provides distributed storage similar to generalized RAID, allowing arbitrary levels of redundancy, unlike RAID's limitations.
// TODO: implement
// Ideally, given S bytes of storage in arbitrary combinations with maximum drive size M, a number R>=M of redundancy bytes can be arbitrarily chosen such that up to R bytes of discrete drives can be lost without losing any data and S-R bytes are usable for storing data. It would probably be an acceptable kludge to group smaller drives into bunches such that each bunch approximates size M, losing some capacity in the process.
//	This should resemble a generalized version of RAID (and similar technologies), overcoming limitations, as follows.
//	* RAID 0 provides no redundancy.
//	* RAID 1 requires S to be an integer multiple of the number of bytes to be stored.
//	* RAID 2-4 are generally considered obsolete.
//	* RAID 5 only supports 1 drive being lost.
//	* RAID 6 only supports up to 2 drives being lost.
//	* BeyondRAID only supports up to 3 drives being lost.

package raidn
