package deck

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"
)

func ExampleCard() {
	fmt.Println(Card{Suit: Spade, Rank: Ace})
	fmt.Println(Card{Suit: Club, Rank: Two})
	fmt.Println(Card{Suit: Diamond, Rank: Queen})
	fmt.Println(Card{Suit: Heart, Rank: Nine})
	fmt.Println(Card{Suit: Joker})
	// Output:
	// Ace of Spades
	// Two of Clubs
	// Queen of Diamonds
	// Nine of Hearts
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 52 {
		t.Errorf("Wrong number of cards in a new deck: got %d, want 52", len(cards))
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	exp := Card{Rank: Ace, Suit: Heart}
	if cards[0] != exp {
		t.Errorf("Expected first card to be Ace of Hearts, got %s", cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	exp := Card{Rank: Ace, Suit: Heart}
	if cards[0] != exp {
		t.Errorf("Expected first card to be Ace of Hearts, got %s", cards[0])
	}
}

func TestJokers(t *testing.T) {
	numJokers := 3
	cards := New(Jokers(numJokers))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != numJokers {
		t.Errorf("Expected %d jokers, got %d", numJokers, count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	cards := New(Filter(filter))
	for _, c := range cards {
		if filter(c) {
			t.Errorf("Did not expect to find %s in deck", c)
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	if len(cards) != 3*52 {
		t.Errorf("Expected %d cards, got %d", 3*52, len(cards))
	}
}

func TestShuffle(t *testing.T) {
	shuffleRand = rand.New(rand.NewSource(0))
	orig := New()
	first := orig[44]
	second := orig[23]
	cards := New(Shuffle)
	if cards[0] != first {
		t.Errorf("Expected first card to be %s, got %s", first, cards[0])
	}
	if cards[1] != second {
		t.Errorf("Expected second card to be %s, got %s", second, cards[1])
	}
}
