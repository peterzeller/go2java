public class Main {

    private static String goFmt(Object o) {
        if (o instanceof java.util.List<?> l) return l.toString().replace(", ", " ");
        return String.valueOf(o);
    }

    public static int fi() {
        return 1;
    }

    public static byte fi8() {
        return 2;
    }

    public static short fi16() {
        return 3;
    }

    public static int fi32() {
        return 4;
    }

    public static long fi64() {
        return 5L;
    }

    public static int fu() {
        return 6;
    }

    public static short fu8() {
        return 255;
    }

    public static short fu8one() {
        return 1;
    }

    public static int fu16() {
        return 65535;
    }

    public static int fu32() {
        return 0x80000000;
    }

    public static long fu64() {
        return 0x8000000000000000L;
    }

    public static long fup() {
        return 9L;
    }

    public static short fb() {
        return 255;
    }

    public static int fr() {
        return 11;
    }

    public static void ops8(short a, short b) {
        System.out.println(goFmt((Integer.compareUnsigned(a, b) <= 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) < 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) > 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) >= 0)));
    }

    public static void ops32(int a, int b) {
        System.out.println(goFmt(Integer.divideUnsigned(a, b)));
        System.out.println(goFmt(Integer.remainderUnsigned(a, b)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) < 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) <= 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) > 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) >= 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) == 0)));
        System.out.println(goFmt((Integer.compareUnsigned(a, b) != 0)));
    }

    public static void ops64(long x, long y) {
        System.out.println(goFmt(Long.divideUnsigned(x, y)));
        System.out.println(goFmt(Long.remainderUnsigned(x, y)));
        System.out.println(goFmt((Long.compareUnsigned(x, y) < 0)));
        System.out.println(goFmt((Long.compareUnsigned(x, y) <= 0)));
        System.out.println(goFmt((Long.compareUnsigned(x, y) > 0)));
        System.out.println(goFmt((Long.compareUnsigned(x, y) >= 0)));
        System.out.println(goFmt((Long.compareUnsigned(x, y) == 0)));
        System.out.println(goFmt((Long.compareUnsigned(x, y) != 0)));
    }

    public static void main(String[] args) {
        System.out.println(goFmt(fi()));
        System.out.println(goFmt(fi8()));
        System.out.println(goFmt(fi16()));
        System.out.println(goFmt(fi32()));
        System.out.println(goFmt(fi64()));
        System.out.println(goFmt(fu()));
        System.out.println(goFmt(fu8()));
        System.out.println(goFmt(fu16()));
        System.out.println(goFmt(fup()));
        System.out.println(goFmt(fb()));
        System.out.println(goFmt(fr()));
        ops8(fu8one(), fu8());
        ops8(fu8(), fu8one());
        ops32(0x80000000, 2);
        ops64(0x8000000000000000L, 3L);
    }
}
