# Note

## Challenge

1. How to calculate flexible size?

   1. How to get size of SomeView?

2. How to set flexible size to views, even ?

## View

### Size

If `initSize` is `-1`, `size` will always be calculated from flexible size.

Otherwise, `size` will be calculated from init size.

## View Modifiers

#### Layout

- Frame
- BackgroundColor
- Padding
- Offset
- Shadow
- Border
- Mask
- Overlay
- Background -> Create New View

#### Action

- OnHover
- OnTapGesture
- OnChange
- OnAppear
- OnDisappear
- OnDrag
- OnDrop

#### Environment

- ForegroundColor
- Font
- FontSize
- FontWeight

#### Parameter

- CornerRadius
- LineSpacing
- Kerning
- Italic

## Size & Position Calculation Flow

#### Concept

```mermaid
sequenceDiagram
   participant P as Parent(Element)
   participant E as Element
   participant C as Children(Element)

   P -->> E: get preload size
   note right of E: <variable> <br> summed size

   loop for each child
      E -->> C: get preload size
      note right of C: recursive

      C -->> E: return size
      E -> E: add size of child into summed size
   end

   E -> E: set size with summed size
   note right of E: size could be <br> constant or flexible(-1)
   E -->> P: return size

   P -->> E: set size + position
   note right of E: <variable> <br> coordinate anchor
   loop for each child
      E -> E: get count of flexible children
      E -> E: calculate flexible size
      E -->> C: set size + position
      note right of C: recursive
   end
   E ->> E: add coordinate anchor with size

```

#### Example

```mermaid
sequenceDiagram
   participant R as Root
   participant V1 as VStack A
   participant V2 as VStack B
   participant V3 as VStack C
   participant V4 as VStack D

   note left of V1: A contains B <br> and B contains C, D

   V3 -> V3: set size with frame
   note right of V3: size could be <br> constant or flexible(-1)

   V3 ->> V2: return size
   V4 -> V4: set size with frame
   note right of V4: size could be <br> constant or flexible(-1)
   V4 ->> V2: return size
   V2 -> V2: set size with frame <br> or summed size of children <br> (frame first)
   note right of V2: size could be <br> constant or flexible(-1)
   V2 ->> V1: return size
   V1 -> V1: set size with frame <br> or summed size of children <br> (frame first)
   note right of V1: size could be <br> constant or flexible(-1)

   R -> R: set window size
   R -> R: get flexible children count
   R -> R: calculate flexible size
   R ->> V1: set unset size with flexible size

   V1 -> V1: get flexible children count
   V1 -> V1: calculate flexible size
   V1 ->> V2: set unset size with flexible size

   V2 -> V2: get flexible children count
   V2 -> V2: calculate flexible size
   V2 ->> V3: set unset size with flexible size
   V2 ->> V3: set position with cache values
   V2 ->> V4: set unset size with flexible size
   V2 ->> V4: set position with cache values
```
