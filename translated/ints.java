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

    public static int fu32one() {
        return 1;
    }

    public static int fu32two() {
        return 2;
    }

    public static int fu32high() {
        return 0x80000000;
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

    public static void loopOps32() {
        var i = 0;
        while ((i < 3)) {
            var a = fu32one();
            if ((i == 1)) {
                        a = fu32two();
                    }
            if ((i == 2)) {
                        a = fu32high();
                    }
            var j = 0;
            while ((j < 3)) {
                        var b = fu32one();
                        if ((j == 1)) {
                                    b = fu32two();
                                }
                        if ((j == 2)) {
                                    b = fu32high();
                                }
                        System.out.println(goFmt((Integer.compareUnsigned(a, b) < 0)));
                        System.out.println(goFmt((Integer.compareUnsigned(a, b) <= 0)));
                        System.out.println(goFmt((Integer.compareUnsigned(a, b) > 0)));
                        System.out.println(goFmt((Integer.compareUnsigned(a, b) >= 0)));
                        System.out.println(goFmt((Integer.compareUnsigned(a, b) == 0)));
                        System.out.println(goFmt((Integer.compareUnsigned(a, b) != 0)));
                        if ((i < 2)) {
                                    if ((j < 2)) {
                                                System.out.println(goFmt(Integer.divideUnsigned(a, b)));
                                                System.out.println(goFmt(Integer.remainderUnsigned(a, b)));
                                            }
                                }
                        j = (j + 1);
                    }
            i = (i + 1);
        }
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
        loopOps32();
    }
}
