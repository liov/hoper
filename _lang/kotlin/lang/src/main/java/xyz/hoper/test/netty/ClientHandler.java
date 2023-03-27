package xyz.hoper.test.netty;


import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerAdapter;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.util.ReferenceCountUtil;
import org.jetbrains.annotations.NotNull;

/**
 * @author ：lbyi
 * @date ：Created in 2019/6/6
 * @description：handle
 */
public class ClientHandler extends ChannelInboundHandlerAdapter {
    @Override
    public void channelActive(@NotNull ChannelHandlerContext ctx) throws Exception {
    }

    @Override
    public void channelRead(@NotNull ChannelHandlerContext ctx, @NotNull Object msg) throws Exception {
        try {
            ByteBuf buf = (ByteBuf) msg;
            byte[] req = new byte[buf.readableBytes()];
            buf.readBytes(req);
            String body = new String(req, "utf-8");
            System.out.println("Client :" + body );
        } finally {
            // 记得释放xxxHandler里面的方法的msg参数: 写(write)数据, msg引用将被自动释放不用手动处理; 但只读数据时,!必须手动释放引用数
            ReferenceCountUtil.release(msg);
        }
    }

    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
            throws Exception {
        ctx.close();
    }

}
