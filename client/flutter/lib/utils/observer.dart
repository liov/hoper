

class Subject<T> {

  List<Observer<T>> observers = List.empty(growable: true);
  T? options;

  void setState(T? options) {
    this.options = options;
    notifyAllObservers();
  }

  void attach(Observer<T> observer){
    observers.add(observer);
  }

  void notifyAllObservers(){
    observers.forEach((observer) { observer.update(options); });
  }
}

abstract class Observer<T> {
    void update(T? options);
}

