/*==========================================================================
module:      Path.cpp
programmer:  Marco Pinter
             Copyright (c) Marco Pinter, 2001.

A* algorithm for path-finding, written from 8/20/2000.

   This code is "open source".
   You may incorporate it royalty-free in any project,
   without restriction.  If possible, please maintain a reference
   to the author, and this message, within the source.

=========================================================================-*/

#include <afxwin.h>
#include <windowsx.h>
#include <math.h>

#define  PATH_CPP
#include "general.h"
#include "linklist.h"
#include "path.h"

#define MAXTILEDIST 40
#define MAXTILE_XLEN (MAXTILEDIST*3/2)
#define MAXTILE_YLEN (MAXTILEDIST*3/2)
#define TILE_CENTER (MAXTILE_XLEN / 2)

#define GetTile(ptsmall) (ptsmall == NULL_TILE ? NULL : &(gPathTileNode[ptsmall.x][ptsmall.y][ptsmall.dir]))



// Rather than using a priority queue, this algorithm speeds things up by:
//  (1) Determining a center point between origin & destination,
//  (2) Declaring a maximum distance from that point which the path can go,
//      and thereby declaring a fixed-sized square in which the path must occur
//  (3) Definining a static array for that square of tiles, which will hold all
//      information normally stored in the A* prioirty queue.

#pragma pack(push, 1)
// Define the structure to be used
struct PathTileNodeType {

    unsigned int  bOpen         : 1;  // Whether this node is in the "Open" list
    unsigned int  direction     : 3;  // Optional: Direction (0-7) at this node
    unsigned int  dummy         : 4;
    unsigned int  costFromStart : 24; // Lowest cost (calculated so far) from origin to
                                      //  this tile
    unsigned int  totalCost     : 24; // Sum of cost to this point + heuristic est (* 64)

    // X,Y coordinates of previous node on path
    unsigned int  parentDir     : 3;
    BYTE          parentX;
    BYTE          parentY;

    // If bOpen, then this tile belongs to an implicit "Open list", &
    //  prevOpen & nextOpen are the list pointers to neighboring elements
    POINTSMALL    prevOpen, nextOpen;

};

// GLOBALS
// Define the static array of tile node info
PathTileNodeType  gPathTileNode[MAXTILE_XLEN][MAXTILE_YLEN][8];

#pragma pack(pop)

// Define array to hold sub-pathing table
SubPathNode       gSubPaths[8][48][8];

// gNextCheck array stores the ordering of a pre-computed list of nodes, for use
//  in Smoothing pass (when bCheckList == true).
short gCheckIndex[MAXTILE_XLEN][MAXTILE_YLEN];
POINT gNextCheck[MIDPOINT_MAX+1];


// external globals
extern bool   gGrid[MAX_GRIDSIZE][MAX_GRIDSIZE];  // The Blocking-tile grid
extern double gTurnRadius;    // turn radius of ship in grid units
extern bool   gbCurved;       // Whether to curve movement, or just use straight lines
extern bool   gbDirectional;  // Whether pathing will formally check directions
extern bool   gbFixed;        // Whether directions are fixed (16 total) or continuous
extern bool   gbSmooth48;     // Whether to smooth the path afterwards: Directional-48
extern bool   gbSmoothSimple; //                                        Simple
extern bool   gbSmooth;       // Whether to smooth the path after it is calculated
extern double gShipWidth;     // width of ship (relative to grid size)
extern bool   gbDirectional;  // Whether pathing will formally check directions
extern double gGranularity;   // How closely the blocking is tested along a path
extern int    gSearch;        // How many surrounding tiles to search
extern bool   gbIgnoreBounds; // Whether we should allow out-of-bounds searches
                              //  (in this test program), to simulate realistic timing

// For speed, we must have a separate array to tell if each node has been visited
// (stores by individual bit)
BYTE gbVisited[MAXTILE_XLEN*MAXTILE_YLEN];

inline void SetVisited(int x, int y, int dir)
{
    int offset = x*MAXTILE_YLEN + y;
    gbVisited[offset] |= (1<<dir);
}

inline bool GetVisited(int x, int y, int dir)
{
    int offset = x*MAXTILE_YLEN + y;
    return ((bool) (gbVisited[offset] & (1<<dir)));
}

// forward declarations
#define CheckGrid(x,y) (x>=0 && y>=0 && x<MAX_GRIDSIZE && y<MAX_GRIDSIZE ? gGrid[x][y] : !gbIgnoreBounds) 

POINTSMALL NULL_TILE = MAKESMALLPOINT(0xFF, 0xFF, 0);
int gSqrtTable[MAXTILE_XLEN][MAXTILE_YLEN];
void SetupSqrtTable()
{
    for (int i=0; i<MAXTILE_XLEN; i++)
        for (int j=0; j<MAXTILE_YLEN; j++)
            gSqrtTable[i][j] = (int) (sqrt(i*i + j*j) * DIST_MULT);
}

POINT gDirOffset[48];
BYTE  gDirIndexFromOffset[7][7];
void SetupDirOffsets()
{
    int xoff, yoff;

    for (int i=0; i<8; i++)
    {
        gDirOffset[i].x = gDirOffset[i].y = 0;
        switch (i)
        {
            case 0:                         gDirOffset[i].y = -1;  break;
            case 1:  gDirOffset[i].x =  1;                         break;
            case 2:                         gDirOffset[i].y =  1;  break;
            case 3:  gDirOffset[i].x = -1;                         break;
            case 4:  gDirOffset[i].x =  1;  gDirOffset[i].y = -1;  break;
            case 5:  gDirOffset[i].x =  1;  gDirOffset[i].y =  1;  break;
            case 6:  gDirOffset[i].x = -1;  gDirOffset[i].y =  1;  break;
            case 7:  gDirOffset[i].x = -1;  gDirOffset[i].y = -1;  break;
        }
    }

    xoff = 2;
    for (i=8; i<=12; i++)
        gDirOffset[i] = MAKEPOINT(xoff, i-10);
    xoff = -2;
    for (i=13; i<=17; i++)
        gDirOffset[i] = MAKEPOINT(xoff, i-15);
    yoff = 2;
    for (i=18; i<=20; i++)
        gDirOffset[i] = MAKEPOINT(i-19, yoff);
    yoff = -2;
    for (i=21; i<=23; i++)
        gDirOffset[i] = MAKEPOINT(i-22, yoff);

    xoff = 3;
    for (i=24; i<=30; i++)
        gDirOffset[i] = MAKEPOINT(xoff, i-27);
    xoff = -3;
    for (i=31; i<=37; i++)
        gDirOffset[i] = MAKEPOINT(xoff, i-34);
    yoff = 3;
    for (i=38; i<=42; i++)
        gDirOffset[i] = MAKEPOINT(i-40, yoff);
    yoff = -3;
    for (i=43; i<=47; i++)
        gDirOffset[i] = MAKEPOINT(i-45, yoff);

    for (i=0; i<48; i++)
        gDirIndexFromOffset[gDirOffset[i].x + 3][gDirOffset[i].y + 3] = i;
}

inline int DistanceMacro(POINT pta, POINT ptb)
{
    return gSqrtTable[(ptb.x > pta.x ? ptb.x - pta.x : pta.x - ptb.x)][(ptb.y > pta.y ? ptb.y - pta.y : pta.y - ptb.y)];
//    return DistanceMacroOld(pta, ptb);
}



int FindPath(FPOINT ptOrigin, FPOINT ptDest, double startDirection, FINFO *pathArray)
{
    // Do the main path-finding routine
    POINT gridOffset;
    int node = DoThePath(ptOrigin, ptDest, startDirection, pathArray, &gridOffset,
                         0xFF, FALSE);
    if (node <= 0)
        return node;

    // Do special smoothing pass
    if (gbSmooth48  &&  gbCurved)
    {
        // Store values of gbDirectional and gSearch
        bool bDirTemp = gbDirectional;
        int  searchTemp = gSearch;

        // Set up for a Directional-48 search
        gbDirectional = true;
        gSearch = 48;

        // Set up the gNextCheck array from the path array
        POINT ptCur = MAKEPOINT((int)ptOrigin.x + gridOffset.x, 
                                (int)ptOrigin.y + gridOffset.y);
        gNextCheck[node] = ptCur;
        gCheckIndex[ptCur.x][ptCur.y] = node;
        for (int i=node-1; i>=0; i--)
        {
            ptCur = MAKEPOINT((int)pathArray[i].pos.x + gridOffset.x, 
                              (int)pathArray[i].pos.y + gridOffset.y);
            gNextCheck[i] = ptCur;
            gCheckIndex[ptCur.x][ptCur.y] = i;
        }

        // Call the DoThePath for special linear search
        node = DoThePath(ptOrigin, ptDest, startDirection, pathArray, &gridOffset,
                         0xFF, true);

        // Restore values
        gbDirectional = bDirTemp;
        gSearch = searchTemp;
    }

    if (node <= 0)
        return node;

    // Do regular smoothing pass
    if (gbSmoothSimple)
        node = SmoothPath(pathArray, node, ptOrigin, startDirection);
    else if (!gbDirectional  &&  !gbSmoothSimple  &&  !gbSmooth48)
        MergeAdjacent(MAKEPOINT((int)ptOrigin.x, (int)ptOrigin.y), pathArray, &node);
    return node;
}


int HybridPath(FPOINT ptOrigin, FPOINT ptDest, double startDirection, FINFO *pathArray)
{
    // Start with the main pathfinding routine (non-directional-4)
    POINT gridOffset;
    gbDirectional = false;
    gSearch = 4;
    int node = DoThePath(ptOrigin, ptDest, startDirection, pathArray, &gridOffset,
                         0xFF, false);
    // if this failed, there is no path
    if (node <= 0)
        return node;

    // Do the first (forward) linear Directional-48 smoothing search
    gbDirectional = true;
    gSearch = 48;
    // Set up the gNextCheck array from the path array
    POINT ptCur = MAKEPOINT((int)ptOrigin.x + gridOffset.x, 
                            (int)ptOrigin.y + gridOffset.y);
    gNextCheck[node] = ptCur;
    gCheckIndex[ptCur.x][ptCur.y] = node;
    for (int i=node-1; i>=0; i--)
    {
        ptCur = MAKEPOINT((int)pathArray[i].pos.x + gridOffset.x, 
                          (int)pathArray[i].pos.y + gridOffset.y);
        gNextCheck[i] = ptCur;
        gCheckIndex[ptCur.x][ptCur.y] = i;
    }
    int origPathLength = node+1;

    // Call the DoThePath for special linear search
    node = DoThePath(ptOrigin, ptDest, startDirection, pathArray, &gridOffset,
                     0xFF, true);

    // If path failed, we need to do more work:
    if (node <=0)
    {
        int   minCost, minDirection, j, nodesA, nodesB, nodesC;
        POINT ptTemp;
        FINFO PathArrayA[MIDPOINT_MAX], PathArrayB[MIDPOINT_MAX], PathArrayC[MIDPOINT_MAX];

        // Find the last node we reached in the search, then find the cheapest final
        //  direction at that node, and build the path to there
        for (i = 0; i < origPathLength; i++)
        {
            ptCur = gNextCheck[i];
            if (gbVisited[ptCur.x*MAXTILE_YLEN + ptCur.y])
            {
                minCost = 9999999;
                for (j=0; j<8; j++)
                    if (GetVisited(ptCur.x, ptCur.y, j)  
                        &&  gPathTileNode[ptCur.x][ptCur.y][j].totalCost < minCost)
                    {
                        minDirection = j;
                        minCost = gPathTileNode[ptCur.x][ptCur.y][j].totalCost;
                    }
                nodesA = ConstructPath(ptOrigin, startDirection, 
                                       MAKEPOINT(ptCur.x - gridOffset.x, ptCur.y - gridOffset.y), 
                                       gridOffset, PathArrayA, minDirection);
                break;
            }
        }

        // Set up gNextSearch for a backwards search from the end
        //  First, compute the direction to launch (from destination) based on
        //   direction of original search
        int endDirection = (gNextCheck[0].x > gNextCheck[1].x ? 4 :
                            (gNextCheck[0].x < gNextCheck[1].x ? 0 :
                             (gNextCheck[0].y < gNextCheck[1].y ? 2 : 6)));

        //  Next, reverse the order of the array
        for (i=0; i<origPathLength / 2; i++)
        {
            ptTemp = gNextCheck[i];
            gNextCheck[i] = gNextCheck[origPathLength - 1 - i];
            gNextCheck[origPathLength - 1 - i] = ptTemp;
        }
        for (i=0; i<origPathLength; i++)
            gCheckIndex[gNextCheck[i].x][gNextCheck[i].y] = i;

        //  Finally, Call the DoThePath for special linear search
        node = DoThePath(ptDest, ptOrigin, Segment2Radian(endDirection), pathArray, &gridOffset, 0xFF, true);

        // Now, as before, find the last node touched before getting stuck
        for (i = 0; i < origPathLength; i++)
        {
            ptCur = gNextCheck[i];
            if (gbVisited[ptCur.x*MAXTILE_YLEN + ptCur.y])
            {
                minCost = 9999999;
                for (j=0; j<8; j++)
                    if (GetVisited(ptCur.x, ptCur.y, j)  
                        &&  gPathTileNode[ptCur.x][ptCur.y][j].totalCost < minCost)
                    {
                        minDirection = j;
                        minCost = gPathTileNode[ptCur.x][ptCur.y][j].totalCost;
                    }
                nodesB = ConstructPath(ptDest, Segment2Radian(endDirection), 
                                       MAKEPOINT(ptCur.x - gridOffset.x, ptCur.y - gridOffset.y), 
                                       gridOffset, PathArrayB, minDirection);
                break;
            }
        }

        // Walk forward in both lists while they're <12 tiles apart
        int ia = 0, ib = 0, lasta = 0, lastb = 0;
        bool bWalkA = true;
        // If they started > 12 tiles apart, fail
        if (fabs(PathArrayA[ia].pos.x - PathArrayB[ib].pos.x) > 12.0  ||
            fabs(PathArrayA[ia].pos.y - PathArrayB[ib].pos.y) > 12.0)
            return -1;

        while (ia < nodesA  ||  ib < nodesB)  
        {
            lasta = ia;
            lastb = ib;
            if (bWalkA  &&  ia < nodesA)
                ia++;
            else if (ib < nodesB)
                ib++;
            bWalkA = !bWalkA;
            if (fabs(PathArrayA[ia].pos.x - PathArrayB[ib].pos.x) > 12.0  ||
                fabs(PathArrayA[ia].pos.y - PathArrayB[ib].pos.y) > 12.0)
            {
                ia = lasta;
                ib = lastb;
                break;
            }
        }


        // Do a special Direction-8 search from the above points
        //  (if we're checking to the destination, it can end on any direction; otherwise fixed)
        gSearch = 8;
        BYTE dirMask = (ib == nodesB ? 0xFF : 1<<((Radian2Segment(PathArrayB[ib].dir)+4)%8) );
        nodesC = DoThePath(PathArrayA[ia].pos, PathArrayB[ib].pos, PathArrayA[ia].dir,
                         PathArrayC, &gridOffset, dirMask, false);
        if (nodesC < 0)
            return nodesC;
        //  (if we were checking to the destination, set the final direction)
        if (ib == nodesB)
                PathArrayB[ib].dir = CapRadian(PathArrayC[0].dir + PI_VAL);

        // Compile the lists
        node = 0;
        //  Put the destination portion of the list on first
        for (i=nodesB; i>=ib; i--)
        {
            pathArray[node] = PathArrayB[i];
            pathArray[node++].dir = CapRadian(PathArrayB[i].dir + PI_VAL);
        }
        //  Next add the special middle section
        for (i=0; i<nodesC; i++)
            pathArray[node++] = PathArrayC[i];
        //  Finally, take the portion of the first list
        for (i=ia; i<nodesA; i++)
            pathArray[node++] = PathArrayA[i];



    }

    // Do regular smoothing pass
    node = SmoothPath(pathArray, node, ptOrigin, startDirection);
    return node;
}




// The main PATHING procedure...
//     dirMaskEnd:     8-bits representing legal ending directions (usually 0xFF)
//     dirMaskStart:   8-bits representing legal start directions (usually 0,
//                      meaning to use ptOrigin and startDirection).  
//                      If != 0, ptStartSpecial is the first location, and we should
//                      push on ALL masked directions as possible starting positions.
//     ptStartSpecial: (See above)
//     bCheckList:     If TRUE, we should only check nodes which are in pre-computed
//                      list, using directions specified in gNextCheck.

// Return values:
//     Positive  # of points in path
//     errors:
//        -1     A path could not be found
//        -2     Points are too far apart
//        -5     Too many midpoints for array size

int  DoThePath(FPOINT ptOrigin, FPOINT ptDest, double startDirection, FINFO *pathArray,
               POINT *pGridOffset, BYTE dirMaskEnd, bool bCheckList)
{
    int   totalDirections = gSearch * (gbDirectional ? 8 : 1);
    int   j, k, xcheck, ycheck, costToGoal, checkBest;
    DWORD dw;
    bool  blocked;
    __int64     posMask;
    SubPathNode subPath;

    // Find grid points for origin & destination
    POINT gdOrigin = { (int)ptOrigin.x, (int)ptOrigin.y };
    POINT gdDest   = { (int)ptDest.x,   (int)ptDest.y   };
    int   dirStart = Radian2Segment(startDirection);

    // Make sure source & dest are not more than MAXTILEDIST tiles apart
    if (abs(gdDest.x - gdOrigin.x) > MAXTILEDIST  ||
        abs(gdDest.y - gdOrigin.y) > MAXTILEDIST)
        return -2;

    // Determine the center tile (between source & dest)
    int xCenter, yCenter;
    xCenter = min(gdOrigin.x, gdDest.x) + abs(gdDest.x - gdOrigin.x) / 2;
    yCenter = min(gdOrigin.y, gdDest.y) + abs(gdDest.y - gdOrigin.y) / 2;

    // For the PathTileNode table, <xCenter,yCenter> will be the center node,
    //  so determine x,y offset to index the table
    POINT gdOffset; 
    gdOffset.x = TILE_CENTER - xCenter;
    gdOffset.y = TILE_CENTER - yCenter;
    *pGridOffset = gdOffset;
    POINT gdDestOff = {gdDest.x + gdOffset.x, gdDest.y + gdOffset.y};
    POINT gdOriginOff = {gdOrigin.x + gdOffset.x, gdOrigin.y + gdOffset.y};

    // Initialize our PathTileNode visited table
    int totalTiles = (MAXTILE_XLEN) * (MAXTILE_YLEN);
    memset(&(gbVisited[0]), 0, totalTiles);

    // Set origin tile data, and "push" onto 'Open list'
    //  (Direction is 0 because later we explicitly calculate origin direction)
    PathTileNodeType *curtile =
               &gPathTileNode[gdOrigin.x + gdOffset.x][gdOrigin.y + gdOffset.y][dirStart];

    curtile->bOpen = true;
    curtile->costFromStart = 0;  
    curtile->totalCost = DistanceMacro(gdOrigin, gdDest);
    curtile->direction = dirStart;
    curtile->parentX = curtile->parentY = 0xFF;
    if (bCheckList)
        // If we are doing a special search thru a pre-computed list, the
        //  heuristic "correct direction" should be 3 tiles ahead in path.
    {
        checkBest = gCheckIndex[gdOriginOff.x][gdOriginOff.y];
        int ii = max(0, checkBest - 2);
        POINT myDest = gNextCheck[ii];
        costToGoal = GetHeuristicDist(DistanceMacro(gdOriginOff, myDest), gdOriginOff, myDest,dirStart)
                     + DIST_MULT * ii;
    }
    else if (gbDirectional)
        curtile->totalCost = GetHeuristicDist(curtile->totalCost, gdOrigin, gdDest, dirStart);
    curtile->nextOpen = curtile->prevOpen = NULL_TILE;
    SetVisited (gdOrigin.x + gdOffset.x, gdOrigin.y + gdOffset.y, dirStart);
    POINTSMALL firstOpen, lastOpen;
    firstOpen = lastOpen = MAKESMALLPOINT(gdOrigin.x + gdOffset.x, gdOrigin.y + gdOffset.y,
                                            dirStart);

    // Loop while there is something on the open list, (or until "broken-out"
    //  due to goal success)
    while (firstOpen != NULL_TILE)
    {
        // Search 'Open list' for the lowest cost node.
        // (This is the slowest part of the algorithm, because it requires a
        //  linear search.)
        PathTileNodeType *searchTile;
        curtile = GetTile(firstOpen);
        for (searchTile = GetTile(curtile->nextOpen); searchTile != NULL;
                                            searchTile = GetTile(searchTile->nextOpen))
            if (searchTile->totalCost < curtile->totalCost)
                curtile = searchTile;

        // Figure out x,y coordinates from the curtile pointer.
        // (Messy, but effective, & better than storing them in the nodes)
        int index = ((BYTE *)curtile - (BYTE *)gPathTileNode) / (sizeof(PathTileNodeType) * 8);
        POINT ptCurrent;
        ptCurrent.x = index / MAXTILE_YLEN;
        ptCurrent.y = index - ptCurrent.x * MAXTILE_YLEN;

        // Check if we've reached the goal
        if (ptCurrent.x == gdDestOff.x  &&  ptCurrent.y == gdDestOff.y  &&
            (dirMaskEnd & (1 << (curtile->direction))))
            return ConstructPath(ptOrigin, startDirection, gdDest, gdOffset, pathArray, curtile->direction);

        // Remove this node from the 'Open list'
        if (curtile->prevOpen != NULL_TILE)
            GetTile(curtile->prevOpen)->nextOpen = curtile->nextOpen;
        else
            firstOpen = curtile->nextOpen;
        if (curtile->nextOpen != NULL_TILE)
            GetTile(curtile->nextOpen)->prevOpen = curtile->prevOpen;
        else
            lastOpen = curtile->prevOpen;
        curtile->bOpen = false;

        // For pre-computed route searches, set flags for legal next positions,
        //  based on 3 ahead (because it's always a Directional-48 search).
        if (bCheckList)
        {
            posMask = 0;
            int ii = gCheckIndex[ptCurrent.x][ptCurrent.y] - 1;
            checkBest = min(ii+1, checkBest);
            while (    ii >= 0  &&  ii <= checkBest + 6
                   &&  abs(ptCurrent.x - gNextCheck[ii].x) <= 3  
                   &&  abs(ptCurrent.y - gNextCheck[ii].y) <= 3)
            {
                posMask |= 1i64 << gDirIndexFromOffset[gNextCheck[ii].x - ptCurrent.x + 3]
                                                      [gNextCheck[ii].y - ptCurrent.y + 3];
                ii--;
            }
        }


        // Loop through all successors of this node (8 surrounding tiles)
        for (int direction = 0; direction < totalDirections; direction ++)
        {
            int dirDest = direction / gSearch;
            int posDest = direction % gSearch;
            // For special pre-computed route searches, check to make sure
            //  next tile is on the route
            if (bCheckList  &&  !(posMask & (1i64<<posDest)))
                continue;

            // Determine coordinates of new tile
            POINT ptNew;
            int costMove;
            ptNew.x = ptCurrent.x + gDirOffset[posDest].x;
            ptNew.y = ptCurrent.y + gDirOffset[posDest].y;

            // Get tile (if valid & not blocked & not origin)
            xcheck = ptNew.x - gdOffset.x;
            ycheck = ptNew.y - gdOffset.y;
            if (ptNew.x < 0  ||  ptNew.x >= MAXTILE_XLEN  ||
                ptNew.y < 0  ||  ptNew.y >= MAXTILE_YLEN  ||
                CheckGrid(xcheck, ycheck)  ||
                (ptNew.x == gdOriginOff.x  &&  ptNew.y == gdOriginOff.y  &&  dirDest == dirStart))
                continue;

            if (!gbDirectional)
                costMove = (posDest < 4 ? (int)BLOCK_NORM : (int)BLOCK_DIAG);
            else
            {
                // Get the sub-path info for this origin and dest
                // (If we're at true origin of path, compute fresh with exact coordinates)
                if (ptCurrent.x == gdOriginOff.x  &&  ptCurrent.y == gdOriginOff.y
                    &&  curtile->direction == dirStart)
                    GetHitsAndLength(&subPath, MAKEFPOINT(5.0 + ptOrigin.x - gdOrigin.x,
                                                          5.0 + ptOrigin.y - gdOrigin.y),
                                     startDirection, MAKEFPOINT(5.5 + gDirOffset[posDest].x, 
                                                                5.5 + gDirOffset[posDest].y),
                                     Segment2Radian(dirDest), gGranularity);
                else
                    subPath = gSubPaths[curtile->direction][posDest][dirDest];

                // Check table against grid values (exit loop if any blockers)
                blocked = false;
                for (j=0; j<4; j++)
                    if (subPath.hitTable[j])
                        for (k=0, dw=1; k<32; k++)
                        {
                            if (subPath.hitTable[j] & dw)
                            {
                                int index = (j<<5) + k;
                                xcheck = ptCurrent.x + index%11 - gdOffset.x - 5;
                                ycheck = ptCurrent.y + index/11 - gdOffset.y - 5;
                                if (CheckGrid(xcheck, ycheck))
                                    blocked = true;
                            }
                            dw = dw<<1;
                        }
                if (blocked)
                    continue;

                costMove = (int)(subPath.length * DIST_MULT);
            }
            PathTileNodeType *newtile = &gPathTileNode[ptNew.x][ptNew.y][dirDest];

            // if node has been visited, & cost to reach it was less, skip it
            int newCost = curtile->costFromStart + costMove;  
            bool  bVisited = GetVisited(ptNew.x, ptNew.y, dirDest);
            if (bVisited  &&  newtile->costFromStart <= newCost) 
                continue;

            // if node hasn't been visited, set member 
            if (!bVisited)
                SetVisited (ptNew.x, ptNew.y, dirDest);

            // Set elements of node & put on 'Open list'
            newtile->parentX = (BYTE) ptCurrent.x;
            newtile->parentY = (BYTE) ptCurrent.y;
            newtile->parentDir = curtile->direction;
            if (bCheckList)
                // If we are doing a special search thru a pre-computed list, the
                //  heuristic "correct direction" should be 3 tiles ahead in path.
            {
                int ii = max(0, gCheckIndex[ptNew.x][ptNew.y] - 2);
                POINT myDest = gNextCheck[ii];
                costToGoal = GetHeuristicDist(DistanceMacro(ptNew, myDest), ptNew, myDest,dirDest)
                             + DIST_MULT * ii;
            }
            else if (gbDirectional)
                // For normal directional searches, compute heuristic as curved path to goal.
                costToGoal = GetHeuristicDist(DistanceMacro(ptNew, gdDestOff), 
                                              ptNew, gdDestOff, dirDest);
            else
                costToGoal = DistanceMacro(ptNew, gdDestOff);
            newtile->costFromStart = newCost;  
            newtile->totalCost = newCost + costToGoal;
            newtile->direction = dirDest;
            POINTSMALL ptNewSmall = MAKESMALLPOINT(ptNew.x, ptNew.y, dirDest);
            if ((!bVisited)  ||  !(newtile->bOpen))
            {
                newtile->bOpen = true;
                newtile->prevOpen = lastOpen;
                newtile->nextOpen = NULL_TILE;
                if (newtile->prevOpen != NULL_TILE)
                    GetTile(newtile->prevOpen)->nextOpen = ptNewSmall;
                lastOpen = ptNewSmall;
                if (firstOpen == NULL_TILE)
                    firstOpen = ptNewSmall;
            }
        } // end for directions
    }

    // If we've reached here, we've failed to find a path.
    // return failure.
    return -1;
}



// Return values:
//     Positive  # of points in path
//     errors:
//        -5     Too many midpoints for array size
int ConstructPath(FPOINT ptOrigin, double startDirection, POINT ptDest, POINT gdOffset, FINFO *pathArray, int finalDir)
{
    int pointCount = 0, curDir;

    // Loop backwards through path, filling in directions
    POINT ptCur = ptDest;
    PathTileNodeType *curtile;
    curtile = &gPathTileNode[ptCur.x + gdOffset.x][ptCur.y + gdOffset.y][finalDir];
    while (curtile->parentX != 0xFF) 
    {
        // Put this node on the array list
        if (pointCount+1 >= MIDPOINT_MAX)
            return -5;
        pathArray[pointCount].pos.x = 0.5 + (double)ptCur.x;
        pathArray[pointCount].pos.y = 0.5 + (double)ptCur.y;
        if (gbDirectional)
            pathArray[pointCount].dir = Segment2Radian(curtile->direction);
        else
            pathArray[pointCount].dir = -1.0;
        pointCount++;

        ptCur.x = curtile->parentX - gdOffset.x;
        ptCur.y = curtile->parentY - gdOffset.y;
        curDir = curtile->parentDir;
        curtile = &gPathTileNode[ptCur.x + gdOffset.x][ptCur.y + gdOffset.y][curDir];
    }
    // Put the origin as an extra node (but don't add to count -- only for use
    //  by hybrid algorithm)
    pathArray[pointCount].pos = ptOrigin;
    pathArray[pointCount].dir = startDirection;

    return pointCount;
}

// For searches which are NOT Directional and will NOT be smoothed, this
//  function simply merges adjacent movements of the same direction, to
//  eliminate extraneous intermediate tiles
void MergeAdjacent(POINT ptOrigin, FINFO *pathArray, int *pNumPoints)
{
    POINT ptLast = MAKEPOINT((int)pathArray[*pNumPoints-1].pos.x, 
                             (int)pathArray[*pNumPoints-1].pos.y);
    POINT ptCurrent;
    int   dxlast = ptLast.x - ptOrigin.x;
    int   dylast = ptLast.y - ptOrigin.y;
    int   dx, dy;
    for (int i=*pNumPoints-1; i>0; i--)
    {
        ptCurrent = MAKEPOINT((int)pathArray[i-1].pos.x, (int)pathArray[i-1].pos.y);
        dx = ptCurrent.x - ptLast.x;
        dy = ptCurrent.y - ptLast.y;
        if (dx == dxlast  &&  dy == dylast)
        {
            for (int j=i; j<*pNumPoints-1; j++)
                pathArray[j].pos = pathArray[j+1].pos;
            (*pNumPoints)--;
        }
        ptLast = ptCurrent;
        dxlast = dx;
        dylast = dy;
    }
}



// Smooth out a computed path, hopefully returning less points in the array.
//  ptOrigin and dirOrigin specify the starting point and direction.
// Return values:
//     Positive  # of points in path
int SmoothPath(FINFO *pathArray, int numPoints, FPOINT ptOrigin, double dirOrigin)
{
    // Define a linked list to be used for points while smoothing
    LinkList<FINFO> moveList;

    // Attempt to smooth from first to last point, and on failure, 
    //  recursively split list and attempt to smooth each half.
    FINFO fiOrigin;
    fiOrigin.pos = ptOrigin;
    fiOrigin.dir = dirOrigin;
//    SmoothSection(pathArray, &moveList, numPoints-1, 0, fiOrigin);
    for (int i=numPoints-1; i>=0; i--)
        moveList.push_back(pathArray[i]);

    // **** Merge adjacent points *****
    // Set up variables:
    LinkList<FINFO>::iterator iter, delnode;
    FINFO curStart = fiOrigin;
    LineSegment lineSeg[4];
    int lastSeg;
    // If curved, the current node always has to have its direction set:
    if (moveList.begin()->data.dir < 0  &&  gbCurved)
    {
        lastSeg = ComputeRoute1(curStart.pos, moveList.begin()->data.pos, 
                                curStart.dir, lineSeg);
        moveList.begin()->data.dir = lineSeg[lastSeg].radStart;
    }
    // The main loop:
    for (iter = moveList.begin(); iter->next != NULL; )
    {
        // Compute the route from origin to destination
        if (!gbCurved)
            MakeStraightLine(curStart.pos, iter->next->data.pos, lineSeg);
        else if (iter->next->data.dir >= 0)
            ComputeRoute2(curStart.pos, iter->next->data.pos, 
                          curStart.dir, iter->next->data.dir, lineSeg);
        else
            lastSeg = ComputeRoute1(curStart.pos, iter->next->data.pos, 
                                    curStart.dir, lineSeg);
        if (gbFixed)
            MakeLineSegmentsFixed(lineSeg);

        // Walk through the route at the gradient level, and see if it's not blocked.
        if (CurveWalk(lineSeg))
        {
            // Advance to the next node, and delete the current one
            delnode = iter;
            iter = iter->next;
            moveList.erase(delnode);
            // If non-directional, fill in the direction, based on the new computed route
            if (iter->data.dir < 0  &&  gbCurved)
                iter->data.dir = lineSeg[lastSeg].radStart;
        }
        else
        {
            // If non-directional, fill in the next direction with the direction that
            //  would result from walking from the previous node.
            // (We don't care if it might be blocked, because the nodes are right next to each
            //  other, and in this mode, we fake it.)
            if (iter->next->data.dir < 0  &&  gbCurved)
            {
                 lastSeg = ComputeRoute1(iter->data.pos, iter->next->data.pos, 
                                         iter->data.dir, lineSeg);
                 iter->next->data.dir = lineSeg[lastSeg].radStart;
            }
            // Advance the start to the current node.  The direction will already be filled in.
            curStart = iter->data;
            // Advance pointer to the next node
            iter = iter->next;
        }
    }
//    MergeAdjacent(&moveList);
//*   (4) Write merge (above) and reconstruct (below)

    // Reconstruct into array
    numPoints = 0;
    for (iter = moveList.begin(); iter != NULL; iter = iter->next)
        pathArray[numPoints++] = iter->data;
    moveList.clear();

    // Reverse order of array
    FINFO finfo;
    for (i=0; i<numPoints/2; i++)
    {
        finfo = pathArray[i];
        pathArray[i] = pathArray[numPoints-1-i];
        pathArray[numPoints-1-i] = finfo;
    }

    return numPoints;
}



bool CurveWalk(LineSegment *lineSeg)
{
    double dist = 0.0;
    double totaldist = lineSeg[0].length + lineSeg[1].length + lineSeg[2].length
                     + lineSeg[3].length;
    FPOINT ptShip;
    double angle;

    // loop thru path, incrementing distance by the gradient at each step
    while (dist < totaldist)
    {
        // Get the point on the path
        CalcVector (lineSeg, dist, &ptShip, &angle);

        // Check the midpoint and four corners to see if they hit blocked tiles
        //  (or out of range)
        int xplus = (int)(ptShip.x + gShipWidth/2),
            yplus = (int)(ptShip.y + gShipWidth/2),
            xminus= (int)(ptShip.x - gShipWidth/2),
            yminus= (int)(ptShip.y - gShipWidth/2);
        if (xminus < 0  ||  yminus < 0  ||  xplus >= MAX_GRIDSIZE  ||  yplus >= MAX_GRIDSIZE)
        {
            if (!gbIgnoreBounds)
                return false;
            if (CheckGrid((int)ptShip.x, (int)ptShip.y)  ||
                CheckGrid(xplus, yplus)  ||
                CheckGrid(xminus, yplus)  ||
                CheckGrid(xplus, yminus)  ||
                CheckGrid(xminus, yminus))
                return false;
        }
        else if (
            gGrid[(int)ptShip.x][(int)ptShip.y]  ||
            gGrid[xplus][yplus]  || 
            gGrid[xminus][yplus]  || 
            gGrid[xplus][yminus]  || 
            gGrid[xminus][yminus]) 

            return false;

        dist += gGranularity;
    }

    return true;
}



void   MakeStraightLine(FPOINT pGridSrc, FPOINT pGridDest, LineSegment *lineSeg)
{
    lineSeg[1].length = lineSeg[2].length = lineSeg[3].length = 0;
    lineSeg[0].bCircle = false;
    lineSeg[0].ptStart = MAKEFPOINT (pGridSrc.x, pGridSrc.y);
    double dx = pGridDest.x - pGridSrc.x;
    double dy = pGridDest.y - pGridSrc.y;
    lineSeg[0].length = sqrt(dx*dx + dy*dy);
    lineSeg[0].radStart = atan2(dy, dx);
    lineSeg[0].radTotal = 0;
}


// Find the best route from Src to Dest, if you have no angle destination
// Pass back the best route in "lineSegment", and
// Return the index of the last actual segment (0-3)
int ComputeRoute1(FPOINT pGridSrc, FPOINT pGridDest,
                   double angleStart, LineSegment *lineSegment)
{
    LineSegment lineSeg[2][4];
    double dlen, minval = 9999999.0;
    int minindex = -1;

    // Try both direct routes.  If one or both possible, use the shortest.
    for (int z=0; z<2; z++)
    {
        bool bClock = (z == 0);
        if (ComputeDirectRoute(pGridSrc, pGridDest, bClock,  angleStart, &(lineSeg[z][0])))
            dlen = lineSeg[z][0].length + lineSeg[z][1].length;
        else
            dlen = 99999999.0;
        if (dlen < minval)
        {
            minval = dlen;
            minindex = z;
        }
    }
    for (int i=0; i<4; i++)
        lineSegment[i] = lineSeg[minindex][i];
    for (i=3; i>= 0; i--)
        if (lineSegment[i].length > 0)
            return i;
}

// Find the best route from Src to Dest, if you have an angle destination
// Pass back the best route in "lineSegment", and
// Return the index of the last actual segment (0-3)
int ComputeRoute2(FPOINT pGridSrc, FPOINT pGridDest,
                   double angleStart, double angleDest, LineSegment *lineSegment)
{
    LineSegment lineSeg[4][4];
    // Try with 2 circles
    double dlen, minval = 9999999.0;
    int minindex = -1;

    for (int z=0; z<4; z++)
    {
        bool bClockA = (z < 2);
        bool bClockB = (z == 0  ||  z == 2);
        if (ComputeLineSegments(pGridSrc, pGridDest, bClockA, bClockB,
                                angleStart, angleDest, &(lineSeg[z][0])))
            dlen = lineSeg[z][0].length + lineSeg[z][1].length + lineSeg[z][2].length;
        else
            dlen = 9999999.0;
        if (dlen < minval)
        {
            minval = dlen;
            minindex = z;
        }
    }
    for (int i=0; i<4; i++)
        lineSegment[i] = lineSeg[minindex][i];
    for (i=3; i>= 0; i--)
        if (lineSegment[i].length > 0)
            return i;
}


// Return the radians point on circle A where the line leaves it to hit
//  circle B.
double FindTouchPoints(FPOINT ptOriginA, bool bClockwiseA, 
                       FPOINT ptOriginB, bool bClockwiseB, double radius,
                       double *pLineLength)
{
    double dx = ptOriginB.x - ptOriginA.x;
    double dy = ptOriginB.y - ptOriginA.y;
    double angleOrigin, angleFinal;

    // determine angle between 2 origin points
    angleOrigin = atan2(dy, dx);

    // determine length of line between origin points
    *pLineLength = sqrt(dx*dx + dy*dy);

    // If A & B rotate in same direction, find touch points by simply using
    //  slope of line, & adding or subtracting 90 degrees as appropriate
    if (bClockwiseA == bClockwiseB)
    {
        return CapRadian(bClockwiseA ? angleOrigin + PI_VAL / 2
                                     : angleOrigin - PI_VAL / 2);
    }
    // If A & B rotate in different directions, more complex calculation
    else if (fabs(2*radius / *pLineLength) <= 1.0)
    {
        // determine angle of right triangle
        double angleRight = acos(2*radius / *pLineLength);

        // add angles for final
        angleFinal = CapRadian(bClockwiseA ? angleOrigin + angleRight 
                                           : angleOrigin - angleRight);

        // determine length of line segment from distance between OriginB,
        //  and the origin of the circle 2R away which is perp. to line
        FPOINT ptNew = { ptOriginA.x + 2*radius*cos(angleFinal),
                         ptOriginA.y + 2*radius*sin(angleFinal) };
        dx = ptOriginB.x - ptNew.x;
        dy = ptOriginB.y - ptNew.y;
        *pLineLength = sqrt(dx*dx + dy*dy);

        return angleFinal;
    }
    else
    {
        *pLineLength = -1.0;
        return 0.0;
    }
}


void MakeLineSegmentsFixed(LineSegment *lineSeg)
{
    double dlen = gTurnRadius / cos(PI_VAL / 16.0);
    double dmax = dlen * sin(PI_VAL / 16.0);
    bool   bClock;
    double length, radStart, radStop, theta, d;
    int    sectionA, sectionB, sdiff;

    // Calculate the actual length of the first circle with discrete segments
    bClock = lineSeg[0].bClockwise;
    radStart= CapRadian(lineSeg[0].radStart);
    radStop = CapRadian(radStart + (bClock ? -lineSeg[0].radTotal : lineSeg[0].radTotal));

    sectionA = ((int)((radStart + PI_VAL / 16.0) / (PI_VAL / 8.0))) % 16;
    sectionB = ((int)((radStop  + PI_VAL / 16.0) / (PI_VAL / 8.0))) % 16;
    sdiff = (bClock ? sectionA - sectionB : sectionB - sectionA);
    if (sdiff < 0)
        sdiff += 16;
    length = (sdiff - 1) * 2.0 * dmax;

    theta = radStart - sectionA * PI_VAL / 8.0;
    d = dlen * sin(theta);
    length += (bClock ? dmax + d : dmax - d);
    
    theta = radStop  - sectionB * PI_VAL / 8.0;
    d = dlen * sin(theta);
    length += (bClock ? dmax - d : dmax + d);
    lineSeg[0].length = length;

    // Calculate the actual length of the 2nd circle (if there is one)
    if (lineSeg[2].length > 0)
    {
        bClock = lineSeg[2].bClockwise;
        radStart= CapRadian(lineSeg[2].radStart);
        radStop = CapRadian(radStart + (bClock ? -lineSeg[2].radTotal : lineSeg[2].radTotal));

        sectionA = ((int)((radStart + PI_VAL / 16.0) / (PI_VAL / 8.0))) % 16;
        sectionB = ((int)((radStop  + PI_VAL / 16.0) / (PI_VAL / 8.0))) % 16;
        sdiff = (bClock ? sectionA - sectionB : sectionB - sectionA);
        if (sdiff < 0)
            sdiff += 16;
        length = (sdiff - 1) * 2.0 * dmax;

        theta = radStart - sectionA * PI_VAL / 8.0;
        d = dlen * sin(theta);
        length += (bClock ? dmax + d : dmax - d);
    
        theta = radStop  - sectionB * PI_VAL / 8.0;
        d = dlen * sin(theta);
        length += (bClock ? dmax - d : dmax + d);
        lineSeg[2].length = length;
    }

    // ***NOTE: The length of the straight line segment(s) is not 100% precise.
    //          The actual length should be reduced slightly, since we are exiting
    //          the first circle and entering the 2nd circle from points slightly
    //          in the middle of them.  If this causes a noticeable jump at the point
    //          between the straight line and the 2nd circle, you can write the
    //          extra code to compute the exact entry & exit points of the line,
    //          and hence the length between them.  (Probably not noticeable, though)

    // Check if the angle is an exact multiple of 22.5 degrees, and
    //  return (done) if it is
    double angle = lineSeg[1].radStart;
    double dmult = (angle * 8.0) / PI_VAL;
    double dclose = (double)((int)(dmult + 0.0001));
    if ((dmult - 0.0001) <= dclose  &&  dclose <= (dmult + 0.0001))
        return;

    // Move the 2nd circle to the 4th spot
    lineSeg[3] = lineSeg[2];

    FPOINT ptDest;
    ptDest.x = lineSeg[1].ptStart.x + lineSeg[1].length * cos(lineSeg[1].radStart);
    ptDest.y = lineSeg[1].ptStart.y + lineSeg[1].length * sin(lineSeg[1].radStart);

    // Divide the straight line into 2 segments, & calculate angle
    //  and length of each
    double m1, m2;  // Calculate the slopes of each path segment
    m1 = tan(dclose * PI_VAL / 8.0);
    m2 = tan((dclose + 1.0) * PI_VAL / 8.0);

    // computing 'b' offset for each line eq.
    double b1, b2;
    b1 = lineSeg[1].ptStart.y - m1 * lineSeg[1].ptStart.x;
    b2 = ptDest.y - m2 * ptDest.x;
    // Handle cases of infinite slope
    lineSeg[2].bCircle = false;
    if (fabs(m1) > 99999)
    {
        lineSeg[2].ptStart.x = lineSeg[1].ptStart.x;
        lineSeg[2].ptStart.y = (m2 * lineSeg[2].ptStart.x) + b2;
    }
    else if (fabs(m2) > 99999)
    {
        lineSeg[2].ptStart.x = ptDest.x;
        lineSeg[2].ptStart.y = (m1 * lineSeg[2].ptStart.x) + b1;
    }
    else
    {

        lineSeg[2].ptStart.x = (b2 - b1) / (m1 - m2);
        lineSeg[2].ptStart.y = (m1 * lineSeg[2].ptStart.x) + b1;
    }

    double dx, dy;
    dx = lineSeg[2].ptStart.x - lineSeg[1].ptStart.x;
    dy = lineSeg[2].ptStart.y - lineSeg[1].ptStart.y;
    lineSeg[1].length = sqrt(dx*dx + dy*dy);
    lineSeg[1].radStart = atan2(dy,dx); 
    
    dx = ptDest.x - lineSeg[2].ptStart.x;
    dy = ptDest.y - lineSeg[2].ptStart.y;
    lineSeg[2].length = sqrt(dx*dx + dy*dy);
    lineSeg[2].radStart = atan2(dy,dx); 

}


// Find a direct route from Src to Dest (if one is possible), starting
//  the circle at angle 'angleStart', in 'bClock' direction.
//  Pass back the filled in LineSegment array.
//  Return TRUE if successful.
bool ComputeDirectRoute(FPOINT pGridSrc, FPOINT pGridDest, bool bClock,
                        double angleStart, LineSegment *lineSeg)
{
    lineSeg[0].bCircle = true;
    lineSeg[0].bClockwise = bClock;

    lineSeg[0].radStart = CapRadian(lineSeg[0].bClockwise ? angleStart + PI_VAL / 2
                                                          : angleStart - PI_VAL / 2);
    lineSeg[0].ptOrigin.x = pGridSrc.x - gTurnRadius * cos(lineSeg[0].radStart);
    lineSeg[0].ptOrigin.y = pGridSrc.y - gTurnRadius * sin(lineSeg[0].radStart);

    // if destination is inside origin circle, we can't get there with a straight line
    double dx, dy;
    dx = pGridDest.x - lineSeg[0].ptOrigin.x;
    dy = pGridDest.y - lineSeg[0].ptOrigin.y;
    double distToOrigin = sqrt(dx*dx + dy*dy);
    if (distToOrigin < gTurnRadius)
        return false;    

    // Figure out length of line from rt angle
    lineSeg[1].length = sqrt(distToOrigin * distToOrigin - gTurnRadius * gTurnRadius);

    // Figure out angle, also from rt triangle
    double angleTriangle = acos (gTurnRadius / distToOrigin);
    double angleFromOrigin = atan2(dy, dx);
    double radStop = (bClock ? angleFromOrigin + angleTriangle : angleFromOrigin - angleTriangle);

    lineSeg[0].radTotal = CapRadian(lineSeg[0].bClockwise ? lineSeg[0].radStart - radStop 
                                                          : radStop - lineSeg[0].radStart);
    lineSeg[0].length = lineSeg[0].radTotal * gTurnRadius;

    // Finish information on the straight line segment (length already set above)
    lineSeg[1].bCircle = false;
    lineSeg[1].ptStart.x = lineSeg[0].ptOrigin.x + gTurnRadius * cos(radStop);
    lineSeg[1].ptStart.y = lineSeg[0].ptOrigin.y + gTurnRadius * sin(radStop);
    lineSeg[1].radStart = CapRadian(lineSeg[0].bClockwise ? radStop - PI_VAL/2 
                                                          : radStop + PI_VAL/2); 
    lineSeg[1].radTotal = 0.0;

    lineSeg[2].length = lineSeg[3].length = 0;
    return true;
}




bool ComputeLineSegments(FPOINT pGridSrc, FPOINT pGridDest, bool bClockA, bool bClockB, 
                         double angleStart, double angleStop, LineSegment *lineSeg)
{
    double radStopA, radStopB;

    lineSeg[0].bCircle = lineSeg[2].bCircle = true;
    lineSeg[0].bClockwise = bClockA;
    lineSeg[2].bClockwise = bClockB;

    lineSeg[0].radStart = CapRadian(lineSeg[0].bClockwise ? angleStart + PI_VAL / 2
                                                          : angleStart - PI_VAL / 2);
    lineSeg[0].ptOrigin.x = pGridSrc.x - gTurnRadius * cos(lineSeg[0].radStart);
    lineSeg[0].ptOrigin.y = pGridSrc.y - gTurnRadius * sin(lineSeg[0].radStart);

    radStopB = CapRadian(lineSeg[2].bClockwise ? angleStop + PI_VAL / 2
                                               : angleStop - PI_VAL / 2);
    lineSeg[2].ptOrigin.x = pGridDest.x - gTurnRadius * cos(radStopB);
    lineSeg[2].ptOrigin.y = pGridDest.y - gTurnRadius * sin(radStopB);

    radStopA = FindTouchPoints(lineSeg[0].ptOrigin, lineSeg[0].bClockwise, 
                               lineSeg[2].ptOrigin, lineSeg[2].bClockwise, 
                               gTurnRadius, &(lineSeg[1].length));
    if (lineSeg[1].length < 0.0)
        return false;

    lineSeg[0].radTotal = CapRadian(lineSeg[0].bClockwise ? lineSeg[0].radStart - radStopA 
                                                          : radStopA - lineSeg[0].radStart);
    lineSeg[0].length = lineSeg[0].radTotal * gTurnRadius;
    lineSeg[2].radStart = (lineSeg[0].bClockwise == lineSeg[2].bClockwise ? 
                                        radStopA : CapRadian(radStopA + PI_VAL));
    lineSeg[2].radTotal = CapRadian(lineSeg[2].bClockwise ? lineSeg[2].radStart - radStopB 
                                                          : radStopB - lineSeg[2].radStart);
    lineSeg[2].length = lineSeg[2].radTotal * gTurnRadius;

    // Finish information on the straight line segment (length already set above)
    lineSeg[1].bCircle = false;
    lineSeg[1].ptStart.x = lineSeg[0].ptOrigin.x + gTurnRadius * cos(radStopA);
    lineSeg[1].ptStart.y = lineSeg[0].ptOrigin.y + gTurnRadius * sin(radStopA);
    lineSeg[1].radStart = CapRadian(lineSeg[0].bClockwise ? radStopA - PI_VAL/2 
                                                          : radStopA + PI_VAL/2); 
    lineSeg[1].radTotal = 0.0;

    lineSeg[3].length = 0;

    return true;
}


// Determine the current point and direction, given the start point
//  and the origin & radians travelled for each of 2 circles used
void CalcVector (LineSegment *lineSeg, double dist, FPOINT *ptNow, double *pSlopeNow)
{
    for (int loop=0; loop<3; loop++)
    {
        LineSegment *ls = &(lineSeg[loop]);
        if (dist <= ls->length  ||  loop == 2)  // don't stop too early on last iteration
        {
            // Is this a circle segment or a straight line segment?
            if (ls->bCircle)
            {
                // Determine the new point on the circle
                double radNew = ls->radStart + (ls->bClockwise ? -1.0 : 1.0) * dist / gTurnRadius;

                // Determine the x,y point from the current angle
                ptNow->x = ls->ptOrigin.x + gTurnRadius * cos(radNew);
                ptNow->y = ls->ptOrigin.y + gTurnRadius * sin(radNew);

                // Determine the new angle (represented in radians, 0.0 points right)
                *pSlopeNow = CapRadian(ls->bClockwise ? radNew - PI_VAL / 2 : radNew + PI_VAL / 2);
                return;
            }
            else  // straight line
            {
                ptNow->x = ls->ptStart.x + dist * cos(ls->radStart);
                ptNow->y = ls->ptStart.y + dist * sin(ls->radStart);
                *pSlopeNow = ls->radStart;
                return;
            }
        }
        dist -= ls->length;
    }

}

    

// Determine the current point and direction for FIXED directions
void CalcVectorFixed (LineSegment *lineSeg, double dist, FPOINT *ptNow, double *pSlopeNow)
{
    for (int loop=0; loop<4; loop++)
    {
        LineSegment *ls = &(lineSeg[loop]);
        if (dist <= ls->length)  
        {
            // Is this a circle segment or a straight line segment?
            if (ls->bCircle)
            {
                double dlen = gTurnRadius / cos(PI_VAL / 16.0);
                double dmax = dlen * sin(PI_VAL / 16.0);
                bool   bClock = ls->bClockwise;
                double radStart= CapRadian(ls->radStart);

                // Find the subsection where this point belongs
                //  First add back to the start of this subsection
                int    section = ((int)((radStart + PI_VAL / 16.0) / (PI_VAL / 8.0))) % 16;
                double theta = radStart - section * PI_VAL / 8.0;
                double d = dlen * sin(theta);
                dist += (bClock ? dmax - d : dmax + d);
                
                // Now count off subsections
                int nsubs = (int) (dist/(2*dmax));
                dist -= nsubs * 2 * dmax;
                section = (bClock ? section - nsubs : section + nsubs);
                double angle = section * (PI_VAL / 8.0);
                *pSlopeNow = angle + (bClock ? -PI_VAL / 2 : PI_VAL / 2);
                FPOINT tpt = { ls->ptOrigin.x + gTurnRadius * cos(angle),
                               ls->ptOrigin.y + gTurnRadius * sin(angle) };
                dist -= dmax;
                ptNow->x = tpt.x + dist * cos(*pSlopeNow);
                ptNow->y = tpt.y + dist * sin(*pSlopeNow);
                return;
            }
            else  // straight line
            {
                ptNow->x = ls->ptStart.x + dist * cos(ls->radStart);
                ptNow->y = ls->ptStart.y + dist * sin(ls->radStart);
                *pSlopeNow = ls->radStart;
                return;
            }
        }
        dist -= ls->length;
    }

}


void SetupHeuristicTable()
{
    FPOINT ptOrigin = {0.0, 0.0}, ptDest = {0.0, 0.0};;
    LineSegment lineSeg[4];
    double dlen1, dlen2, dlen, angle;

    for (int j=0; j<8; j++)
        gHeuristic[0][j] = 0;

    for (int i=1; i<=MAX_HDIST_LONG; i++)
    {
        for (int j=0; j<8; j++)
        {
            angle = (double)j * (PI_VAL / 8.0);
            ptDest.x = ((double) i) / DIST_MULT;
            dlen1 = dlen2 = 99999999.0;
            if (ComputeDirectRoute(ptOrigin, ptDest, true,  angle, lineSeg))
                dlen1 = lineSeg[0].length + lineSeg[1].length;
            if (ComputeDirectRoute(ptOrigin, ptDest, false,  angle, lineSeg))
                dlen2 = lineSeg[0].length + lineSeg[1].length;
            dlen = min(dlen1, dlen2);

            gHeuristic[i][j] = (int)(DIST_MULT * dlen1);
        }
    }
}
    

void GetHitsAndLength(SubPathNode *node, FPOINT ptOrigin, double dirOrigin,
                      FPOINT ptDest, double dirDest, double granularity)
{
    double dist = 0.0;
    memset(&(node->hitTable[0]), 0, 16);
    LineSegment lineSeg[4];
    FPOINT ptShip;
    double angle;
    int index;

    // Compute the route and its distance
    ComputeRoute2(ptOrigin, ptDest, dirOrigin, dirDest, lineSeg);
    node->length = (float) (lineSeg[0].length + lineSeg[1].length + lineSeg[2].length
                            + lineSeg[3].length);
    if (gbFixed)
        MakeLineSegmentsFixed(lineSeg);

    // loop thru path, incrementing distance by the gradient at each step
    while (dist < node->length)
    {
        // Get the point on the path
        CalcVector (lineSeg, dist, &ptShip, &angle);

        // Check the midpoint and four corners to see where they hit
        int xplus = (int)(ptShip.x + gShipWidth/2),
            yplus = (int)(ptShip.y + gShipWidth/2),
            xminus= (int)(ptShip.x - gShipWidth/2),
            yminus= (int)(ptShip.y - gShipWidth/2);

        index = (int)ptShip.y * 11 + (int)ptShip.x;
        if (0 <= index  &&  index <= 111)
            node->hitTable[index>>5] |= 1<<(index%32);

        index = (int)yplus * 11 + (int)xplus;
        if (0 <= index  &&  index <= 111)
            node->hitTable[index>>5] |= 1<<(index%32);

        index = (int)yplus * 11 + (int)xminus;
        if (0 <= index  &&  index <= 111)
            node->hitTable[index>>5] |= 1<<(index%32);

        index = (int)yminus * 11 + (int)xplus;
        if (0 <= index  &&  index <= 111)
            node->hitTable[index>>5] |= 1<<(index%32);

        index = (int)yminus * 11 + (int)xminus;
        if (0 <= index  &&  index <= 111)
            node->hitTable[index>>5] |= 1<<(index%32);

        dist += granularity;
    }

    // Finally, set the origin & destination nodes back to 0, because
    //  we don't really need to check them
    index = (int)ptOrigin.y * 11 + (int)ptOrigin.x;
    node->hitTable[index>>5] &= ((1<<(index%32)) ^ 0xFFFFFFFF);

    index = (int)ptDest.y * 11 + (int)ptDest.x;
    node->hitTable[index>>5] &= ((1<<(index%32)) ^ 0xFFFFFFFF);
}


void SetupTurnTable()
{
    for (int dir1 = 0; dir1 < 8; dir1++)
        for (int endpos = 0; endpos < 48; endpos++)
            for (int dir2 = 0; dir2 < 8; dir2++)
            {
                GetHitsAndLength(&(gSubPaths[dir1][endpos][dir2]), MAKEFPOINT(5.5,5.5),
                    Segment2Radian(dir1), 
                    MAKEFPOINT(5.5 + gDirOffset[endpos].x, 5.5 + gDirOffset[endpos].y),
                    Segment2Radian(dir2), 0.1);
            }
}
