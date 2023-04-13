package lorem

// just a small helper to produce Lorem-ish text

import (
	"math/rand"
	"strings"
	"time"
)

const latin = "At vero eos et accusamus et iusto odio dignissimos ducimus, qui blanditiis praesentium voluptatum deleniti atque corrupti, quos dolores et quas molestias excepturi sint, obcaecati cupiditate non provident, similique sunt in culpa, qui officia deserunt mollitia animi, id est laborum et dolorum fuga. Et harum quidem rerum facilis est et expedita distinctio. Nam libero tempore, cum soluta nobis est eligendi optio, cumque nihil impedit, quo minus id, quod maxime placeat, facere possimus, omnis voluptas assumenda est, omnis dolor repellendus. Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet, ut et voluptates repudiandae sint et molestiae non recusandae. Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat"

func GenerateLorem(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var text string
	words := strings.Fields(latin)
	for i := 0; i < length; i++ {
		word := words[r.Intn(len(words))]
		text += word + " "
	}
	return text
}
