# gotils-xmap

Fast iteration golang map.

# Use example use case

Classic pub sub scenario. xmap will be a topic and will hold references to each subcribers to this topic. Iterate through xmap to push data to each subscriber. FILO ordering because if iterated in the usual way.

## Pros

    - Fast iteration
    - preserves ordering
    - Fast lookup via map

## Cons

    - Slow insert functions
    - watch for size - shard the topics if neccessary
