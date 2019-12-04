#ifndef _HASHTABLE_H
#define _HASHTABLE_H
#include <stdint.h>
#include <stdlib.h>

/*****************************************************************************
 * This library defines a generic hashtable. Both the key and data can be any
 * data type. Each is defined by a pointer to a location in memory and a size
 * of the memory structure in bytes. The hashtable always copies the key into
 * internal memory which it will free when destroyed. The data can be copied
 * into memory as well (and is with the hashtable_insert command) or a
 * reference to memory managed by the user will be maintained, using the 
 * hashtable_insertref command.
 *
 * hashtable_create:
 * This command creates a hashtable of size size with an optional hashfunction
 * that you can define. If you choose NULL for the hashfunction, it will 
 * default to a 64 bit FNV-1a hash. The hash function should return a 64 bit
 * unsigned integer. The size is merely the size of the hash array. More
 * elements can be added, in the event of a hash collision the second hash is
 * added to the same hash bucket. The hashtable_create function returns a
 * pointer to a hashtable or NULL if it fails.
 *
 * hashtable_insert:
 * This will insert new data or overwrite old data (with the same key). It 
 * copies both the key and data into newly allocated memory. The key and data
 * can be anything. The function returns <0 on error.
 *
 * hashtable_insertref:
 * This is like hashtable_insert except that it only adds a pointer to the
 * data structure, and does not copy the data itself.
 *
 * hashtable_remove:
 * Removes an item from the hashtable and frees the relevant memory. Returns
 * <0 if unsuccessful.
 *
 * hashtable_get:
 * Returns a pointer to the data associated with the key or NULL if the value
 * was not found.
 *
 * hashtable_free:
 * Frees all memory associated with the hashtable. Remember to use this.
 * **************************************************************************/

struct hashnode {
  void *key;
  size_t key_size;
  void *data;
  size_t data_size;
  struct hashnode *next;
};

typedef struct _hashtable {
  struct hashnode **nodearray;
  uint64_t (*hashfunc)(void *, size_t);
  size_t size;
} hashtable;

typedef struct _hashtable_iterator {
  hashtable *iterating_table;
  size_t table_row;
  struct hashnode *current_node;
} hashtable_iterator;

hashtable *hashtable_create(size_t size, uint64_t (*hashfunc)(void *, size_t));
int hashtable_insert(hashtable *hashtbl, void *key, size_t key_size, void *data, size_t data_size);
int hashtable_insertref(hashtable *hashtbl, void *key, size_t key_size, void *data);
int hashtable_remove(hashtable *hashtbl, void *key, size_t key_size);
void *hashtable_get(hashtable *hashtbl, void *key, size_t key_size);
void hashtable_free(hashtable *hashtbl);
hashtable_iterator *hashtable_iterator_create(hashtable *hashtable);
void hashtable_iterator_next(hashtable_iterator *iterator);
void *hashtable_iterator_get_key(hashtable_iterator *iterator);
size_t hashtable_iterator_get_key_size(hashtable_iterator *iterator);
void *hashtable_iterator_get_data(hashtable_iterator *iterator);
size_t hashtable_iterator_get_data_size(hashtable_iterator *iterator);
void hashtable_iterator_free(hashtable_iterator *iterator);

uint64_t fnv1a64(void *buf, size_t len);

#endif /*!_HASHTABLE_H*/
