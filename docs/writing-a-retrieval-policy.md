# Writing a retrieval policy

`airecall.yaml` controls how memories are retrieved. Getting this right
matters more than the model choice — a bad policy makes the agent forget
what matters and remember what doesn't.

## Modes

| Mode | What runs | Best for |
|---|---|---|
| `keyword` | BM25-style term overlap | exact terms, product names, ids |
| `vector` | cosine over word-overlap vectors | paraphrase, synonyms |
| `hybrid` | 0.65 keyword + 0.35 vector | default; most cases |

## Rules of thumb

- **Keep `top_k` small** (3-8). Memory injected into the prompt costs
  tokens and dilutes attention; the top 3 relevant episodes beat 20

