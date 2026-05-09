public record Vec2(int x, int y) {
    public Vec2 plus(Vec2 other) {
        return new Vec2((x() + other.x), (y() + other.y));
    }
    public static void main(String[] args) {
        var a = new Vec2(1, 2);
        var b = new Vec2(3, 4);
        var c = a.plus(b);
        System.out.println(c.x);
        System.out.println(c.y);
    }
}
