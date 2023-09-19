package HashTable

// https://leetcode.cn/problems/two-sum
func twoSum(nums []int, target int) []int {
	numToIndexMap := make(map[int]int, len(nums))
	for i, num := range nums {
		if index, ok := numToIndexMap[target-num]; ok {
			return []int{index, i}
		}
		numToIndexMap[num] = i
	}
	return []int{}
}