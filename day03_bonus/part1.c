#include "hashtable.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct _point {
  int x;
  int y;
} point;

void build_grid(hashtable *grid, char *path);

int main(int argc, char *argv[]) {
  FILE *fp;
  char *path_one = NULL;
  size_t len_one;
  size_t len_two;
  char *path_two = NULL;
  char *token;
  hashtable *grid_one;
  hashtable *grid_two;
  hashtable_iterator *iterator;
  point *position_key;
  int distance;
  int i;
  int min_distance;
  size_t pt_size = sizeof(point);

  fp = fopen("input.txt", "r");
  if(fp == NULL) {
    fprintf(stderr, "Unable to open file\n");
    exit(1);
  }

  if(getline(&path_one, &len_one, fp)<0) {
    fprintf(stderr, "Unable to read first path\n");
    exit(1);
  }

  if(getline(&path_two, &len_two, fp)<0) {
    fprintf(stderr, "Unable to read second path\n");
    exit(1);
  }
  fclose(fp);

  grid_one = hashtable_create(65536, NULL);
  if(grid_one == NULL) {
    fprintf(stderr, "Unable to create first grid\n");
    exit(1);
  }
  build_grid(grid_one, path_one);

  grid_two = hashtable_create(65536, NULL);
  if(grid_two == NULL) {
    fprintf(stderr, "Unable to create second grid\n");
    exit(1);
  }
  build_grid(grid_two, path_two);

  iterator = hashtable_iterator_create(grid_one);
  min_distance = 0;
  hashtable_iterator_next(iterator);
  while((position_key=(point*)hashtable_iterator_get_key(iterator))!=NULL) {
    if(hashtable_get(grid_two, position_key, pt_size)!=NULL) {
      distance = abs(position_key->x)+abs(position_key->y);
      if(min_distance==0 || distance<min_distance) {
        min_distance=distance;
      }
    }
    hashtable_iterator_next(iterator);
  }
  hashtable_iterator_free(iterator);
  hashtable_free(grid_one);
  hashtable_free(grid_two);
  printf("Shortest distance: %d\n", min_distance);
}

void build_grid(hashtable *grid, char *path) {
  point position;
  size_t pt_size = sizeof(position);
  char wire = 1;
  char *token = NULL;
  char direction;
  int distance;
  int i;

  position.x = 0;
  position.y = 0;
  while((token = strsep(&path, ",")) != NULL) {
    direction = token[0];
    distance = atoi(&token[1]);
    for(i=0; i<distance; i++) {
      switch(direction) {
        case 'U':
          position.y++;
          break;
        case 'D':
          position.y--;
          break;
        case 'R':
          position.x++;
          break;
        case 'L':
          position.x--;
          break;
        default:
          fprintf(stderr, "Unknown direction: %c\n", direction);
          exit(1);
      }

      hashtable_insert(grid, &position, pt_size, &wire, 1);
    }
  }
}
