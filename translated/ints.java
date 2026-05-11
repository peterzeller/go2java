public class Main {

    private static String goFmt(Object o) {
        if (o instanceof java.util.List<?> l) return l.toString().replace(", ", " ");
        return String.valueOf(o);
    }

    public static java.util.List<Integer> ints() {
        return new java.util.ArrayList<>(java.util.List.of((-1), 1, 2147483647));
    }

    public static java.util.List<Byte> int8s() {
        return new java.util.ArrayList<>(java.util.List.of((byte)(-1), (byte)1, (byte)127));
    }

    public static java.util.List<Short> int16s() {
        return new java.util.ArrayList<>(java.util.List.of((short)(-1), (short)1, (short)32767));
    }

    public static java.util.List<Integer> int32s() {
        return new java.util.ArrayList<>(java.util.List.of((-1), 1, 2147483647));
    }

    public static java.util.List<Long> int64s() {
        return new java.util.ArrayList<>(java.util.List.of((long)(-1), (long)1L, (long)9223372036854775807L));
    }

    public static java.util.List<Integer> uints() {
        return new java.util.ArrayList<>(java.util.List.of(1, 2, 2147483647));
    }

    public static java.util.List<Short> uint8s() {
        return new java.util.ArrayList<>(java.util.List.of((short)1, (short)128, (short)255));
    }

    public static java.util.List<Integer> uint16s() {
        return new java.util.ArrayList<>(java.util.List.of(1, 32768, 65535));
    }

    public static java.util.List<Integer> uint32s() {
        return new java.util.ArrayList<>(java.util.List.of(1, 2, 0x80000000));
    }

    public static java.util.List<Long> uint64s() {
        return new java.util.ArrayList<>(java.util.List.of((long)1L, (long)2L, (long)0x8000000000000000L));
    }

    public static java.util.List<Long> uintptrs() {
        return new java.util.ArrayList<>(java.util.List.of((long)1L, (long)2L, (long)9L));
    }

    public static java.util.List<Short> bytesVals() {
        return new java.util.ArrayList<>(java.util.List.of((short)1, (short)128, (short)255));
    }

    public static java.util.List<Integer> runesVals() {
        return new java.util.ArrayList<>(java.util.List.of((-1), 1, 1114111));
    }

    public static void main(String[] args) {
        System.out.println(goFmt(ints().get(0)));
        System.out.println(goFmt(int8s().get(1)));
        System.out.println(goFmt(int16s().get(2)));
        System.out.println(goFmt(int32s().get(0)));
        System.out.println(goFmt(int64s().get(1)));
        System.out.println(goFmt(uints().get(2)));
        System.out.println(goFmt(uint8s().get(2)));
        System.out.println(goFmt(uint16s().get(1)));
        System.out.println(goFmt(uint32s().get(1)));
        System.out.println(goFmt(uint64s().get(1)));
        System.out.println(goFmt(uintptrs().get(2)));
        System.out.println(goFmt(bytesVals().get(2)));
        System.out.println(goFmt(runesVals().get(2)));
        var vals = uint32s();
        var i = 0;
        while ((i < 3)) {
            var j = 0;
            while ((j < 3)) {
                        var a = vals.get(i);
                        var b = vals.get(j);
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
}
