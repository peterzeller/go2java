public class Vec2 {
    private int X;
    private int Y;

    public Vec2(int X, int Y) {
        this.X = X;
        this.Y = Y;
    }

    public int getX() { return this.X; }
    public void setX(int X) { this.X = X; }

    public int getY() { return this.Y; }
    public void setY(int Y) { this.Y = Y; }

    public Vec2 plus(Vec2 other) {
        return new Vec2((this.X + other.getX()), (this.Y + other.getY()));
    }

    public void moveX(int delta) {
        this.X = (this.X + delta);
    }

    public static void main(String[] args) {
        var a = new Vec2(1, 2);
        var b = new Vec2(3, 4);
        var c = a.plus(b);
        c.moveX(10);
        System.out.println(c.getX());
        System.out.println(c.getY());
    }
}
