/***********************************************************************/
/***********************************************************************/
/***                                                                 ***/
/***    Path.h:  Include file for pathing algorithm                  ***/
/***                                                                 ***/
/***                                                                 ***/
/***********************************************************************/
/***********************************************************************/

#define MIDPOINT_MAX   100
#define MAX_GRIDSIZE   30

#define DIST_MULT      64

const unsigned int BLOCK_NORM =  DIST_MULT;
const unsigned int BLOCK_DIAG = (int)(1.4142 * DIST_MULT);

#define CLOSE_VAL  0.00000001
const double PI_VAL_2 = 2.0 * PI_VAL;
const double PI_VAL_2CLOSE = PI_VAL_2 - CLOSE_VAL;

class POINTSMALL
{
public:
    BYTE x;
    BYTE y;
    BYTE dir;

    inline POINTSMALL &operator =(const POINTSMALL &cval)   {x=cval.x; y=cval.y; dir = cval.dir; return *this;}
//    inline POINTSMALL &operator =(const POINT &cval)   {x=(BYTE)cval.x; y=(BYTE)cval.y; return *this;}
    inline bool		operator ==(const POINTSMALL &cval)  {return (x==cval.x && y==cval.y  &&  dir==cval.dir);}
    inline bool		operator !=(const POINTSMALL &cval)  {return !(this->operator==(cval));}

};

inline POINTSMALL MAKESMALLPOINT(int xval, int yval, int dir)
{
    POINTSMALL p;
    p.x = (BYTE) xval;
    p.y = (BYTE) yval;
    p.dir = (BYTE) dir;
    return p;
}

struct FINFO
{
    FPOINT  pos;
    double  dir;
};


struct LineSegment
{
    bool    bCircle;    // Is it a circle curve - T (or a straight line - F)

    FPOINT  ptOrigin;   // For circles, the origin of the circle
    bool    bClockwise; // For circles, if the initial direction increases radians
    FPOINT  ptStart;    // For straight lines, the starting point of this line

    double  radStart;   // The radians point at start of curve (or slope for straight line)
    double  radTotal;   // The total radians of the curve (if one)
    double  length;     // The total length of this line segment
};


struct SubPathNode
{
    DWORD   hitTable[4];    // 121 bit table (9x9) showing which tiles get hit
                            //  when going from origin(4,4,d1) to dest(x,y,d2)
    float   length;         // Length of path going from origin to dest
};


// Constrain a radian value to 0<=value<=2*PI
inline double CapRadian(double r)
{
    while (r >= 2.0*PI_VAL)
        r -= 2.0*PI_VAL;
    while (r < 0.0)
        r += 2.0*PI_VAL;
    if (r < CLOSE_VAL  ||  r > PI_VAL_2CLOSE)
        return 0.0;
    return r;
}

// Macros for changing between angle segments (0-7) and true radian values
#define Radian2Segment(rad) ((int)(CapRadian(rad + PI_VAL / 8.0) / (PI_VAL / 4.0)))
#define Segment2Radian(seg) (((double)seg) * PI_VAL / 4.0)

// Calculate the difference between two angles
inline double AngleDiff (double angle1, double angle2)
{
    // Calculate from both 0 and 180 degrees, and take minimum
    double calc1, calc2;
    calc1 = fabs(angle1-angle2);

    if (angle1 > PI_VAL)
        angle1 -= 2*PI_VAL;
    if (angle2 > PI_VAL)
        angle2 -= 2*PI_VAL;
    calc2 = fabs(angle1-angle2);
    return min(calc1, calc2);
}

// Heuristic definitions & calculations
#ifdef PATH_CPP

#define MAX_HDIST         10
#define MAX_HDIST_LONG    MAX_HDIST*DIST_MULT

short gHeuristic[MAX_HDIST_LONG+1][8];

inline short GetHeuristicDist(int dist, POINT ptOrigin, POINT ptDest, int dirToDest)
{
    double angleCorrect = CapRadian(atan2(ptDest.y - ptOrigin.y, ptDest.x - ptOrigin.x));
    int dirdiff = (int) (AngleDiff(angleCorrect, Segment2Radian(dirToDest)) / (PI_VAL / 8.0));
    dirdiff = min(7, dirdiff);

    if (dist <= MAX_HDIST_LONG)
        return gHeuristic[dist][dirdiff];
    // if the points are more than 10 tiles apart, approximate distance
    //  using the heuristic for 10 tiles
    else
        return gHeuristic[MAX_HDIST_LONG][dirdiff]
                + (dist - MAX_HDIST_LONG);
}
#endif


int    FindPath(FPOINT ptOrigin, FPOINT ptDest, double startDirection, FINFO *pathArray);
int    HybridPath(FPOINT ptOrigin, FPOINT ptDest, double startDirection, FINFO *pathArray);

void   MergeAdjacent(POINT ptOrigin, FINFO *pathArray, int *pNumPoints);

void   MakeStraightLine(FPOINT pGridSrc, FPOINT pGridDest, LineSegment *lineSeg);

int    SmoothPath(FINFO *pathArray, int numPoints, FPOINT ptOrigin, double dirOrigin);

int    ComputeRoute1(FPOINT pGridSrc, FPOINT pGridDest,
                     double angleStart, LineSegment *lineSegment);
int    ComputeRoute2(FPOINT pGridSrc, FPOINT pGridDest,
                     double angleStart, double angleDest, LineSegment *lineSegment);

void SetupSqrtTable();
void SetupDirOffsets();
void SetupHeuristicTable();
void SetupTurnTable();


// Internal pathing procedures
int    DoThePath(FPOINT ptOrigin, FPOINT ptDest, double startDirection, FINFO *pathArray,
                 POINT *pGridOffset, BYTE dirMaskEnd, bool bCheckList);
int    ConstructPath(FPOINT ptOrigin, double startDirection, POINT ptDest, POINT gdOffset, FINFO *pathArray, int finalDir);
bool   ComputeDirectRoute(FPOINT pGridSrc, FPOINT pGridDest, bool bClock,
                          double angleStart, LineSegment *lineSeg);

void   SmoothSection(FINFO *pathArray, LinkList<FINFO> *moveList, 
                     int ifirst, int ilast, FINFO ptOrigin);
bool   CurveWalk(LineSegment *lineSeg);
double FindTouchPoints(FPOINT ptOriginA, bool bClockwiseA, 
                       FPOINT ptOriginB, bool bClockwiseB, double radius,
                       double *pLineLength);
bool   ComputeLineSegments(FPOINT pGridSrc, FPOINT pGridDest, bool bClockA, bool bClockB, 
                         double angleStart, double angleStop, LineSegment *lineSeg);
void   MakeLineSegmentsFixed(LineSegment *lineSeg);
void   CalcVector (LineSegment *lineSeg, double dist, FPOINT *ptNow, double *pSlopeNow);
void   CalcVectorFixed (LineSegment *lineSeg, double dist, FPOINT *ptNow, double *pSlopeNow);
void   GetHitsAndLength(SubPathNode *node, FPOINT ptOrigin, double dirOrigin,
                        FPOINT ptDest, double dirDest, double granularity);


