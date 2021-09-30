class BeeKeeper {
    hasMask: boolean;
}

class ZooKeeper {
    nametag: string;
}

class GAnimal {
    numLegs: number;
}

class Bee extends GAnimal {
    keeper: BeeKeeper;
}

class Lion extends GAnimal {
    keeper: ZooKeeper;
}

function createInstance<A extends GAnimal>(c: new () => A): A {
    return new c();
}

createInstance(Lion).keeper.nametag;  // typechecks!
createInstance(Bee).keeper.hasMask;   // typechecks!
