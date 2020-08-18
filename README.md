# gotils-xmap

Fast iteration golang map.

# Use example use case

Classic pub sub scenario. xmap will be a topic and will hold references to each subcribers to this topic. Iterate through xmap to push data to each subscriber. FILO ordering because if iterated in the usual way.

Arose because of the need to quickly iterate through sets and retrieve objects quickly. The insert and delete functions can be traded for an increase in core requirement spec.

## CRUD options

    - Add ```
    xm := XMap()

    // Add function - keys must be strings
    xm.Add(key, value)

    // Eco Add function
    // Looks for free spaces in the xmap structure and fills them
    // returns the index of the newly inserted item
    newIndex := AddEco(k, v)
    ```

    - Delete
    -```
    // Delete function
    xm.Delete(key)

    // Delete by item index
    xm.DeleteElementAt(itemIndex)
    ```

    - Update
    ```
    // Overwrite with the Insert function with the same key value
    // (like a normal map)
    xm.Add(key, newValue)

    // Update value via the item index
    // returns the previous and current values
    prev, current := xm.Update(itemIndex, newValue)
    ```

    - Read
    ```
    // Read an item by a key
    // returns a pointer to the value and a Boolean of success
    pointerToValue, ok := xm.GetByKey(key)


    // Read by index
    pointerToValue := xm.GetByIndex(itemIndex)
    ```

    - Memory management
    ```
    // Reindex the structure by squashing all free space
    // This will also reset the indicies for all items
    // Only use when you no longer need previous index
    xm.RebuildIndex()

    // Concurrently calculates free space within the structure
    // Gives manual control of the process
    xm.IndexFreeSpace()

    // View free space
    // returns an int of the number of free slots in the structure
    count := xm.FreeSpace()
    ```

### Add is via

## Pros

    - Fast iteration
    - preserves ordering
    - Fast lookup via map
    - Key <=> index lookup
    - Preserves the map type's value versatility Accepts any / and mixed datatypes as values in the same structure

## Cons

    - Slow insert functions
    - Size can be an issue for very large sets since the underlying xmap uses a slice to preserve the order
