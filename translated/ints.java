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
        return 5;
    }

    public static int fu() {
        return 6;
    }

    public static byte fu8() {
        return 100;
    }

    public static short fu16() {
        return 1000;
    }

    public static int fu32() {
        return 7;
    }

    public static long fu64() {
        return 8;
    }

    public static long fup() {
        return 9;
    }

    public static byte fb() {
        return 10;
    }

    public static int fr() {
        return 11;
    }

    public static void ops32(int a, int b) {
        System.out.println(goFmt(Integer.divideUnsigned(a, b)));
        System.out.println(goFmt(Integer.remainderUnsigned(a, b)));
        System.out.println(goFmt((a < b)));
        System.out.println(goFmt((a <= b)));
        System.out.println(goFmt((a > b)));
        System.out.println(goFmt((a >= b)));
        System.out.println(goFmt((a == b)));
        System.out.println(goFmt((a != b)));
        System.out.println(goFmt((a + b)));
        System.out.println(goFmt((a - b)));
        System.out.println(goFmt((a * b)));
    }

    public static void ops64(long x, long y) {
        System.out.println(goFmt(Long.divideUnsigned(x, y)));
        System.out.println(goFmt(Long.remainderUnsigned(x, y)));
        System.out.println(goFmt((x < y)));
        System.out.println(goFmt((x <= y)));
        System.out.println(goFmt((x > y)));
        System.out.println(goFmt((x >= y)));
        System.out.println(goFmt((x == y)));
        System.out.println(goFmt((x != y)));
        System.out.println(goFmt((x + y)));
        System.out.println(goFmt((x - y)));
        System.out.println(goFmt((x * y)));
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
        System.out.println(goFmt(fu32()));
        System.out.println(goFmt(fu64()));
        System.out.println(goFmt(fup()));
        System.out.println(goFmt(fb()));
        System.out.println(goFmt(fr()));
        ops32(7, 2);
        ops64(7, 3);
    }
}
