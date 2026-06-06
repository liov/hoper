

class Subject<T> {
  Subject(this.options);
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
    for (var observer in observers) { observer.update(options); }
  }
}

abstract class Observer<T> {
    void update(T? options);
}

