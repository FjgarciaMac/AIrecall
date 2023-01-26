# Writing a retrieval policy

`airecall.yaml` controls how memories are retrieved. Getting this right
matters more than the model choice — a bad policy makes the agent forget
what matters and remember what doesn't.

## Modes

| Mode | What runs | Best for |
|---|---|---|
| `keyword` | BM25-style term overlap | exact terms, product names, ids |
| `vector` | cosine over word-overlap vectors | paraphrase, synonyms |
