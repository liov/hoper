using System;
namespace RectangleApplication
{
    class Rectangle
    {
        // 成员变量
        double length;
        double width;
        public void Acceptdetails()
        {
            length = 4.5;    
            width = 3.5;
        }
        public double GetArea()
        {
            return length * width;
        }
        public void Display()
        {
            Console.WriteLine("Length: {0}", length);
            Console.WriteLine("Width: {0}", width);
            Console.WriteLine("Area: {0}", GetArea());
        }
    }
    
    class ExecuteRectangle
    {
        static void Main(string[] args)
        {
            Rectangle r = new Rectangle();
            r.Acceptdetails();
            r.Display();
            Console.ReadLine();
        }
        
        [TestMethod]
        public void subTest(){
            Stopwatch sw = new Stopwatch(); 
            sw.Start();
            sleep(1000);
            sw.Stop(); 
            TimeSpan ts2 = sw.Elapsed; 
            Console.WriteLine("Stopwatch总共花费{0}ms.", ts2.TotalMilliseconds); 
        }
    }
}
