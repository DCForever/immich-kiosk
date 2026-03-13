# Collage Layout Diagrams

Visual diagrams for each collage grid layout. Letters (a, b, c, …) map to `nth-child` indices 1, 2, 3, ….

---

## 1 photo

### collage-1

```
+-----+
|  a  |
+-----+
```

---

## 2 photos

### collage-2

```
+-----+-----+
|     |     |
|  a  |  b  |
|     |     |
+-----+-----+
```
Cells span 2 rows (1.2fr 1fr).

---

## 3 photos

### collage-3 / collage-3-a

```
+----+----+
|    | b  |
| a  +----+
|    | c  |
+----+----+
```
Hero left (a spans 2 rows).

### collage-3-b

```
+----+----+
| a  |    |
+----+ b  |
| c  |    |
+----+----+
```
Hero right (b spans 2 rows).

### collage-3-c

```
+---------+
|    a    |
+----+----+
| b  | c  |
+----+----+
```
Hero top (a spans 2 cols).

---

## 4 photos

### collage-4 / collage-4-a

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
```

### collage-4-b

```
+----+----+
|    a    |
+----+----+
| b  | c  |
+----+----+
```
Hero top (a spans 2 cols).

### collage-4-c

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
|       d       |
+----+----+----+
```
Hero bottom (d spans 3 cols).

---

## 5 photos

### collage-5 / collage-5-a

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
```
Hero left (a spans 2 rows).

### collage-5-b

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
|    e    |
+---------+
```
Hero bottom (e spans 2 cols).

### collage-5-c

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  |    e    |
+----+---------+
```
Hero bottom-right (e spans 2 cols).

---

## 6 photos

### collage-6 / collage-6-a

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
| e  | f  |
+----+----+
```
Top row taller (1.2fr).

### collage-6-b

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
```
Middle row taller (1.2fr).

### collage-6-c

```
+----+----+
|    | b  |
| a  +----+
|    | c  |
+----+----+
| d  | e  |
+----+----+
```
Hero left (a spans 2 rows).

---

## 7 photos

### collage-7 / collage-7-a

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
|    | f  | g  |
+----+----+----+
```
Hero left (a spans 3 rows).

### collage-7-b

```
+----+----+----+----+
|    a    | b  | c  |
+----+----+----+----+
| d  | e  | f  | g  |
|    |    |    |    |
+----+----+----+----+
```
Hero top (a spans 2 cols); d,e,f,g span 2 rows.

### collage-7-c

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
| e  | f  |
+----+----+
|    g    |
+---------+
```
Hero bottom (g spans 2 cols).

---

## 8 photos

### collage-8 / collage-8-a

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
| f  | g  | h  |
+----+----+----+
```
Hero left (a spans 2 rows).

### collage-8-b

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
```
First col wider (1.2fr); middle row taller (1.2fr).

### collage-8-c

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
| e  | f  |
+----+----+
| g  | h  |
+----+----+
```
First col wider (1.2fr); top row taller (1.2fr).

---

## 9 photos

### collage-9 / collage-9-a

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
```
First col wider (1.2fr); middle row taller (1.2fr).

### collage-9-b

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
| f  | g  | h  |
+----+----+----+
|       i       |
+---------------+
```
Hero left (a spans 2 rows); hero bottom (i spans 3 cols).

### collage-9-c

```
+----+----+----+
|    |    | b  |
| a  | a  +----+
|    |    | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
```
Hero top-left (a spans 2×2).

---

## 10 photos

### collage-10 / collage-10-a

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  |
+----+----+
```
First col wider (1.2fr); middle row taller (1.2fr).

### collage-10-b

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
|       j       |
+---------------+
```
First col wider (1.2fr); third row taller (1.2fr); j spans 3 cols.

### collage-10-c

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
| f  | g  | h  |
+----+----+----+
| i  |    j    |
+----+---------+
```
Hero left (a spans 2 rows); j spans 2 cols.

---

## 11 photos

### collage-11 / collage-11-a

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  | k  | .  |
+----+----+----+----+
```
First col wider (1.2fr); middle row taller (1.2fr).

### collage-11-b

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
| j  |    k    |
+----+---------+
```
k spans 2 cols.

### collage-11-c

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
| f  | g  | h  |
+----+----+----+
| i  | j  | k  |
+----+----+----+
```
Hero left (a spans 2 rows).

---

## 12 photos

### collage-12 / collage-12-a

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  | k  | l  |
+----+----+----+----+
```

### collage-12-b

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
| j  | k  | l  |
+----+----+----+
```
First col wider (1.2fr); third row taller (1.2fr).

### collage-12-c

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
| e  | f  |
+----+----+
| g  | h  |
+----+----+
| i  | j  |
+----+----+
| k  | l  |
+----+----+
```
First col wider (1.2fr); second row taller (1.2fr).

---

## 13 photos

### collage-13 / collage-13-a

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  | k  | l  |
+----+----+----+----+
| m  | .  | .  | .  |
+----+----+----+----+
```
First col wider (1.2fr); third row taller (1.2fr).

### collage-13-b

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
| j  | k  | l  |
+----+----+----+
|       m       |
+---------------+
```
m spans 3 cols.

### collage-13-c

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
| f  | g  | h  |
+----+----+----+
| i  | j  | k  |
+----+----+----+
| l  |    m    |
+----+---------+
```
Hero left (a spans 2 rows); m spans 2 cols.

---

## 14 photos

### collage-14 / collage-14-a

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  | k  | l  |
+----+----+----+----+
| m  | n  | .  | .  |
+----+----+----+----+
```
First col wider (1.2fr); second row taller (1.2fr).

### collage-14-b

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
| j  | k  | l  |
+----+----+----+
| m  |    n    |
+----+---------+
```
n spans 2 cols.

### collage-14-c

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
| e  | f  |
+----+----+
| g  | h  |
+----+----+
| i  | j  |
+----+----+
| k  | l  |
+----+----+
| m  | n  |
+----+----+
```
First col wider (1.2fr); second row taller (1.2fr).

---

## 15 photos

### collage-15 / collage-15-a

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  | k  | l  |
+----+----+----+----+
| m  | n  | o  | .  |
+----+----+----+----+
```
First col wider (1.2fr); second row taller (1.2fr).

### collage-15-b

```
+----+----+----+
| a  | b  | c  |
+----+----+----+
| d  | e  | f  |
+----+----+----+
| g  | h  | i  |
+----+----+----+
| j  | k  | l  |
+----+----+----+
| m  | n  | o  |
+----+----+----+
```
First col wider (1.2fr); third row taller (1.2fr).

### collage-15-c

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  | k  | l  |
+----+----+----+----+
| m  | n  |    o    |
+----+----+---------+
```
o spans 2 cols.

---

## 16 photos

### collage-16 / collage-16-a

```
+----+----+----+----+
| a  | b  | c  | d  |
+----+----+----+----+
| e  | f  | g  | h  |
+----+----+----+----+
| i  | j  | k  | l  |
+----+----+----+----+
| m  | n  | o  | p  |
+----+----+----+----+
```
First col wider (1.2fr); second row taller (1.2fr).

### collage-16-b

```
+----+----+
| a  | b  |
+----+----+
| c  | d  |
+----+----+
| e  | f  |
+----+----+
| g  | h  |
+----+----+
| i  | j  |
+----+----+
| k  | l  |
+----+----+
| m  | n  |
+----+----+
| o  | p  |
+----+----+
```
First col wider (1.2fr); second row taller (1.2fr).

### collage-16-c

```
+----+----+----+
|    | b  | c  |
| a  +----+----+
|    | d  | e  |
+----+----+----+
| f  | g  | h  |
+----+----+----+
| i  | j  | k  |
+----+----+----+
| l  | m  | n  |
+----+----+----+
| o  |    p    |
+----+---------+
```
Hero left (a spans 2 rows); p spans 2 cols.
