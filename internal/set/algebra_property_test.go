package set

import (
	"maps"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSet_AlgebraAgainstMapOracle(t *testing.T) {
	random := rand.New(rand.NewSource(42))

	for caseIdx := range 200 {
		baseValues := randomInts(random, random.Intn(12))
		operandValues := [][]int{
			randomInts(random, random.Intn(12)),
			randomInts(random, random.Intn(12)),
			randomInts(random, random.Intn(12)),
		}
		operands := []Source[int]{
			sliceCollection[int](operandValues[0]),
			NewHashSetFromSlice(operandValues[1]),
			sliceCollection[int](operandValues[2]),
		}

		base := intMembership(baseValues)
		operandMemberships := []map[int]struct{}{
			intSourceMembership(operands[0]),
			intSourceMembership(operands[1]),
			intSourceMembership(operands[2]),
		}
		union := unionMembership(base, operandMemberships)
		intersection := intersectionMembership(base, operandMemberships)
		difference := differenceMembership(base, operandMemberships)
		symmetricDifference := symmetricDifferenceMembership(base, operandMemberships)
		relationIdx := caseIdx % len(operands)
		relationOperand := operands[relationIdx]
		relationMembership := operandMemberships[relationIdx]

		original := NewHashSetFromSlice(baseValues)
		assert.Equal(t, union, hashSetMembership(original.Union(operands...)), caseIdx)
		assert.Equal(t, intersection, hashSetMembership(original.Intersection(operands...)), caseIdx)
		assert.Equal(t, difference, hashSetMembership(original.Difference(operands...)), caseIdx)
		assert.Equal(t, symmetricDifference, hashSetMembership(
			original.SymmetricDifference(operands...),
		), caseIdx)
		assert.Equal(t, base, hashSetMembership(original), caseIdx)

		assert.Equal(t, maps.Equal(base, relationMembership), original.Equal(relationOperand), caseIdx)
		assert.Equal(t, isSubsetMembership(base, relationMembership), original.IsSubset(relationOperand), caseIdx)
		assert.Equal(t, isProperSubsetMembership(base, relationMembership), original.IsProperSubset(relationOperand), caseIdx)
		assert.Equal(t, isSubsetMembership(relationMembership, base), original.IsSuperset(relationOperand), caseIdx)
		assert.Equal(t, isProperSubsetMembership(relationMembership, base), original.IsProperSuperset(relationOperand), caseIdx)
		assert.Equal(t, isDisjointMembership(base, relationMembership), original.IsDisjoint(relationOperand), caseIdx)

		mutable := original.Clone()
		assert.Equal(t, len(union)-len(base), mutable.UnionWith(operands...), caseIdx)
		assert.Equal(t, union, hashSetMembership(mutable), caseIdx)

		mutable = original.Clone()
		assert.Equal(t, len(base)-len(intersection), mutable.IntersectWith(operands...), caseIdx)
		assert.Equal(t, intersection, hashSetMembership(mutable), caseIdx)

		mutable = original.Clone()
		assert.Equal(t, len(base)-len(difference), mutable.DifferenceWith(operands...), caseIdx)
		assert.Equal(t, difference, hashSetMembership(mutable), caseIdx)

		mutable = original.Clone()
		assert.Equal(t, membershipChanges(base, symmetricDifference), mutable.SymmetricDifferenceWith(operands...), caseIdx)
		assert.Equal(t, symmetricDifference, hashSetMembership(mutable), caseIdx)
	}
}

func TestKeyedHashSet_AlgebraAgainstMapOracle(t *testing.T) {
	random := rand.New(rand.NewSource(84))

	for caseIdx := range 200 {
		baseValues := randomRecords(random, random.Intn(12))
		operandValues := [][]record{
			randomRecords(random, random.Intn(12)),
			randomRecords(random, random.Intn(12)),
			randomRecords(random, random.Intn(12)),
		}
		operands := []Source[record]{
			sliceCollection[record](operandValues[0]),
			NewKeyedHashSetFromSlice(operandValues[1], recordNameLength),
			NewKeyedHashSetFromSlice(operandValues[2], recordID),
		}

		base := recordMembership(NewKeyedHashSetFromSlice(baseValues, recordID))
		operandMemberships := []map[int]struct{}{
			recordMembership(operands[0]),
			recordMembership(operands[1]),
			recordMembership(operands[2]),
		}
		union := unionMembership(base, operandMemberships)
		intersection := intersectionMembership(base, operandMemberships)
		difference := differenceMembership(base, operandMemberships)
		symmetricDifference := symmetricDifferenceMembership(base, operandMemberships)
		relationIdx := caseIdx % len(operands)
		relationOperand := operands[relationIdx]
		relationMembership := operandMemberships[relationIdx]

		original := NewKeyedHashSetFromSlice(baseValues, recordID)
		assert.Equal(t, union, recordMembership(original.Union(operands...)), caseIdx)
		assert.Equal(t, intersection, recordMembership(original.Intersection(operands...)), caseIdx)
		assert.Equal(t, difference, recordMembership(original.Difference(operands...)), caseIdx)
		assert.Equal(t, symmetricDifference, recordMembership(
			original.SymmetricDifference(operands...),
		), caseIdx)
		assert.Equal(t, base, recordMembership(original), caseIdx)

		assert.Equal(t, maps.Equal(base, relationMembership), original.Equal(relationOperand), caseIdx)
		assert.Equal(t, isSubsetMembership(base, relationMembership), original.IsSubset(relationOperand), caseIdx)
		assert.Equal(t, isProperSubsetMembership(base, relationMembership), original.IsProperSubset(relationOperand), caseIdx)
		assert.Equal(t, isSubsetMembership(relationMembership, base), original.IsSuperset(relationOperand), caseIdx)
		assert.Equal(t, isProperSubsetMembership(relationMembership, base), original.IsProperSuperset(relationOperand), caseIdx)
		assert.Equal(t, isDisjointMembership(base, relationMembership), original.IsDisjoint(relationOperand), caseIdx)

		mutable := original.Clone()
		assert.Equal(t, len(union)-len(base), mutable.UnionWith(operands...), caseIdx)
		assert.Equal(t, union, recordMembership(mutable), caseIdx)

		mutable = original.Clone()
		assert.Equal(t, len(base)-len(intersection), mutable.IntersectWith(operands...), caseIdx)
		assert.Equal(t, intersection, recordMembership(mutable), caseIdx)

		mutable = original.Clone()
		assert.Equal(t, len(base)-len(difference), mutable.DifferenceWith(operands...), caseIdx)
		assert.Equal(t, difference, recordMembership(mutable), caseIdx)

		mutable = original.Clone()
		assert.Equal(t, membershipChanges(base, symmetricDifference), mutable.SymmetricDifferenceWith(operands...), caseIdx)
		assert.Equal(t, symmetricDifference, recordMembership(mutable), caseIdx)
	}
}

func randomInts(random *rand.Rand, length int) []int {
	values := make([]int, length)
	for idx := range values {
		values[idx] = random.Intn(9) - 4
	}
	return values
}

func randomRecords(random *rand.Rand, length int) []record {
	values := make([]record, length)
	for idx := range values {
		values[idx] = record{
			ID:   random.Intn(9) - 4,
			Name: strings.Repeat("x", random.Intn(5)),
		}
	}
	return values
}

func recordNameLength(value record) int {
	return len(value.Name)
}

func intMembership(values []int) map[int]struct{} {
	membership := make(map[int]struct{}, len(values))
	for _, value := range values {
		membership[value] = struct{}{}
	}
	return membership
}

func intSourceMembership(source Source[int]) map[int]struct{} {
	membership := make(map[int]struct{}, source.Length())
	for value := range source.All() {
		membership[value] = struct{}{}
	}
	return membership
}

func hashSetMembership(source *HashSet[int]) map[int]struct{} {
	return maps.Clone(source.values)
}

func recordMembership(source Source[record]) map[int]struct{} {
	membership := make(map[int]struct{}, source.Length())
	for value := range source.All() {
		membership[value.ID] = struct{}{}
	}
	return membership
}

func unionMembership(
	base map[int]struct{},
	operands []map[int]struct{},
) map[int]struct{} {
	result := maps.Clone(base)
	for _, operand := range operands {
		for value := range operand {
			result[value] = struct{}{}
		}
	}
	return result
}

func intersectionMembership(
	base map[int]struct{},
	operands []map[int]struct{},
) map[int]struct{} {
	result := maps.Clone(base)
	for _, operand := range operands {
		for value := range result {
			if _, exists := operand[value]; !exists {
				delete(result, value)
			}
		}
	}
	return result
}

func differenceMembership(
	base map[int]struct{},
	operands []map[int]struct{},
) map[int]struct{} {
	result := maps.Clone(base)
	for _, operand := range operands {
		for value := range operand {
			delete(result, value)
		}
	}
	return result
}

func symmetricDifferenceMembership(
	base map[int]struct{},
	operands []map[int]struct{},
) map[int]struct{} {
	result := maps.Clone(base)
	for _, operand := range operands {
		for value := range operand {
			if _, exists := result[value]; exists {
				delete(result, value)
			} else {
				result[value] = struct{}{}
			}
		}
	}
	return result
}

func isSubsetMembership(
	subset, superset map[int]struct{},
) bool {
	for value := range subset {
		if _, exists := superset[value]; !exists {
			return false
		}
	}
	return true
}

func isProperSubsetMembership(
	subset, superset map[int]struct{},
) bool {
	return len(subset) < len(superset) && isSubsetMembership(subset, superset)
}

func isDisjointMembership(
	left, right map[int]struct{},
) bool {
	for value := range left {
		if _, exists := right[value]; exists {
			return false
		}
	}
	return true
}

func membershipChanges(
	left, right map[int]struct{},
) int {
	changed := 0
	for value := range left {
		if _, exists := right[value]; !exists {
			changed++
		}
	}
	for value := range right {
		if _, exists := left[value]; !exists {
			changed++
		}
	}
	return changed
}
