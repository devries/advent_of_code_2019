#include "hashtable.h"
#include <inttypes.h>
#include <stdio.h>
#include <string.h>

const uint64_t fnv64_prime = UINT64_C(1099511628211);
const uint64_t fnv64_offset = UINT64_C(14695981039346656037);

uint64_t fnv1a64(void *buf, size_t len) {
  uint8_t *pointer = (uint8_t *)buf;
  uint8_t *buf_end = pointer+len;
  uint64_t hash = fnv64_offset;

  while(pointer < buf_end) {
    hash ^= (uint64_t)*pointer++;
    hash *= fnv64_prime;
  }

  return hash;
}

hashtable *hashtable_create(size_t size, uint64_t (*hashfunc)(void *, size_t)) {
  hashtable *hashtbl;
  int i;

  hashtbl = (hashtable *)malloc(sizeof(hashtable));
  if(!hashtbl) {
    return NULL;
  }

  hashtbl->nodearray = (struct hashnode**)malloc(size*sizeof(struct hashnode*));
  if(!hashtbl->nodearray) {
    free(hashtbl);
    return NULL;
  }

  for(i=0;i<size;i++) {
    hashtbl->nodearray[i]=NULL;
  }
  
  hashtbl->size=size;
  if(hashfunc) {
    hashtbl->hashfunc=hashfunc;
  }
  else {
    hashtbl->hashfunc=fnv1a64;
  }

  return hashtbl;
}

int hashtable_insert(hashtable *hashtbl, void *key, size_t key_size, void *data, size_t data_size) {
  struct hashnode *curr_node;
  size_t hashpos;
  size_t olddata_size;
  void *olddata;

  hashpos = hashtbl->hashfunc(key,key_size)%hashtbl->size;

  curr_node = hashtbl->nodearray[hashpos];
  while(curr_node) {
    if(key_size==curr_node->key_size) {
      if(memcmp(key,curr_node->key,key_size)==0) {
        olddata = curr_node->data;
        olddata_size = curr_node->data_size;

        curr_node->data_size=data_size;
        curr_node->data = malloc(data_size);
        
        if(!curr_node->data) {
          curr_node->data_size=olddata_size;
          curr_node->data=olddata;
          return -1;
        }

        if(olddata_size!=0) {
          free(olddata);
        }

        memcpy(curr_node->data,data,data_size);
        return 0;
      }
    }
    curr_node=curr_node->next;
  }

  curr_node = (struct hashnode*)malloc(sizeof(struct hashnode));
  if(!curr_node) {
    return -1;
  }


  curr_node->key_size=key_size;
  curr_node->key = malloc(key_size);
  if(!curr_node->key) {
    free(curr_node);
    return -1;
  }
  memcpy(curr_node->key,key,key_size);

  curr_node->data_size=data_size;
  curr_node->data = malloc(data_size);
  if(!curr_node->data) {
    free(curr_node->key);
    free(curr_node);
    return -1;
  }
  memcpy(curr_node->data,data,data_size);
  
  curr_node->next = hashtbl->nodearray[hashpos];
  hashtbl->nodearray[hashpos]=curr_node;

  return 0;
}

int hashtable_insertref(hashtable *hashtbl, void *key, size_t key_size, void *data) {
  struct hashnode *curr_node;
  size_t hashpos;

  hashpos = hashtbl->hashfunc(key,key_size)%hashtbl->size;

  curr_node = hashtbl->nodearray[hashpos];
  while(curr_node) {
    if(key_size==curr_node->key_size) {
      if(memcmp(key,curr_node->key,key_size)==0) {
        if(curr_node->data_size!=0) {
          free(curr_node->data);
        }

        curr_node->data_size=0;
        curr_node->data=data;
        return 0;
      }
    }
    curr_node=curr_node->next;
  }

  curr_node = (struct hashnode*)malloc(sizeof(struct hashnode));
  if(!curr_node) {
    return -1;
  }

  curr_node->key_size=key_size;
  curr_node->key = malloc(key_size);
  if(!curr_node->key) {
    free(curr_node);
    return -1;
  }
  memcpy(curr_node->key,key,key_size);

  curr_node->data_size=0;
  curr_node->data = data;
  
  curr_node->next = hashtbl->nodearray[hashpos];
  hashtbl->nodearray[hashpos]=curr_node;

  return 0;
}

int hashtable_remove(hashtable *hashtbl, void *key, size_t key_size) {
  struct hashnode *curr_node;
  size_t hashpos;
  struct hashnode *prev_node;
  
  hashpos = hashtbl->hashfunc(key,key_size)%hashtbl->size;

  curr_node = hashtbl->nodearray[hashpos];
  prev_node = NULL;

  while(curr_node) {
    if(key_size==curr_node->key_size) {
      if(memcmp(key,curr_node->key,key_size)==0) {
        free(curr_node->key);
        if(curr_node->data_size!=0) {
          free(curr_node->data);
        }
        if(prev_node) {
          prev_node->next=curr_node->next;
        }
        else {
          hashtbl->nodearray[hashpos]=curr_node->next;
        }
        free(curr_node);
        return 0;
      }
    }
    prev_node=curr_node;
    curr_node=curr_node->next;
  }
  return -1;
}

void *hashtable_get(hashtable *hashtbl, void *key, size_t key_size) {
  struct hashnode *curr_node;
  size_t hashpos;
  
  hashpos = hashtbl->hashfunc(key,key_size)%hashtbl->size;

  curr_node = hashtbl->nodearray[hashpos];

  while(curr_node) {
    if(key_size==curr_node->key_size) {
      if(memcmp(key,curr_node->key,key_size)==0) {
        return curr_node->data;
      }
    }
    curr_node=curr_node->next;
  }
  return NULL;
}

void hashtable_free(hashtable *hashtbl) {
  int i;
  struct hashnode *curr_node;
  struct hashnode *next_node;

  for(i=0;i<hashtbl->size;i++) {
    curr_node = hashtbl->nodearray[i];
    while(curr_node) {
      next_node=curr_node->next;
      free(curr_node->key);
      if(curr_node->data_size!=0) {
        free(curr_node->data);
      }
      free(curr_node);
      curr_node=next_node;
    }
  }

  free(hashtbl->nodearray);
  free(hashtbl);
}


hashtable_iterator *hashtable_iterator_create(hashtable *hashtable) {
  hashtable_iterator *iterator;

  iterator = (hashtable_iterator *)malloc(sizeof(hashtable_iterator));
  if(!iterator) {
    return NULL;
  }

  iterator->iterating_table = hashtable;
  iterator->table_row = -1;
  iterator->current_node = NULL;

  return iterator;
}

void hashtable_iterator_next(hashtable_iterator *iterator) {
  if(iterator->current_node!=NULL) {
    iterator->current_node=iterator->current_node->next;
  }

  while(iterator->current_node==NULL) {
    iterator->table_row++;
    if(iterator->table_row>=iterator->iterating_table->size) break;
    iterator->current_node = iterator->iterating_table->nodearray[iterator->table_row];
  }
}

void *hashtable_iterator_get_key(hashtable_iterator *iterator) {
  if(iterator->current_node==NULL) {
    return NULL;
  }

  return iterator->current_node->key;
}

size_t hashtable_iterator_get_key_size(hashtable_iterator *iterator) {
  if(iterator->current_node==NULL) {
    return 0;
  }
  
  return iterator->current_node->key_size;
}

void *hashtable_iterator_get_data(hashtable_iterator *iterator) {
  if(iterator->current_node==NULL) {
    return NULL;
  }

  return iterator->current_node->data;
}

size_t hashtable_iterator_get_data_size(hashtable_iterator *iterator) {
  if(iterator->current_node==NULL) {
    return 0;
  }

  return iterator->current_node->data_size;
}

void hashtable_iterator_free(hashtable_iterator *iterator) {
  free(iterator);
}

