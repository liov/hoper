import 'package:app/global/service.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/cupertino.dart';

class HeroImage extends StatelessWidget{
  const  HeroImage({super.key,
    required this.child,
    required this.tag,
    required this.slidePageKey,
    this.slideType = SlideType.onlyImage,
  });

  final Widget child;
  final SlideType slideType;
  final Object tag;
  final GlobalKey<ExtendedImageSlidePageState> slidePageKey;

  @override
  Widget build(BuildContext context) {
    late final RectTween rectTween;
    return Hero(
      tag: tag,
      createRectTween: (Rect? begin, Rect? end) {
         rectTween = RectTween(begin: begin, end: end);
        globalService.logger.fine(rectTween);
        return rectTween;
      },
      // make hero better when slide out
      flightShuttleBuilder: (BuildContext flightContext,
          Animation<double> animation,
          HeroFlightDirection flightDirection,
          BuildContext fromHeroContext,
          BuildContext toHeroContext) {
        // make hero more smoothly
        final Hero hero = (flightDirection == HeroFlightDirection.pop
            ? fromHeroContext.widget
            : toHeroContext.widget) as Hero;
        if (flightDirection == HeroFlightDirection.pop) {
          final bool fixTransform = slideType == SlideType.onlyImage &&
              (slidePageKey.currentState!.offset != Offset.zero ||
                  slidePageKey.currentState!.scale != 1.0);

          final Widget toHeroWidget = (toHeroContext.widget as Hero).child;
          return AnimatedBuilder(
            animation: animation,
            builder: (BuildContext buildContext, Widget? child) {
              Widget animatedBuilderChild = hero.child;

              // make hero more smoothly
              animatedBuilderChild = Stack(
                clipBehavior: Clip.antiAlias,
                alignment: Alignment.center,
                children: <Widget>[
                  Opacity(
                    opacity: 1 - animation.value,
                    child: UnconstrainedBox(
                      child: SizedBox(
                        width: rectTween.begin!.width,
                        height: rectTween.begin!.height,
                        child: toHeroWidget,
                      ),
                    ),
                  ),
                  Opacity(
                    opacity: animation.value,
                    child: animatedBuilderChild,
                  )
                ],
              );

              // fix transform when slide out
              if (fixTransform) {
                final Tween<Offset> offsetTween = Tween<Offset>(
                    begin: Offset.zero,
                    end: slidePageKey.currentState!.offset);

                final Tween<double> scaleTween = Tween<double>(
                    begin: 1.0, end: slidePageKey.currentState!.scale);
                animatedBuilderChild = Transform.translate(
                  offset: offsetTween.evaluate(animation),
                  child: Transform.scale(
                    scale: scaleTween.evaluate(animation),
                    child: animatedBuilderChild,
                  ),
                );
              }

              return animatedBuilderChild;
            },
          );
        }
        return hero.child;
      },
      child: child,
    );
  }
}