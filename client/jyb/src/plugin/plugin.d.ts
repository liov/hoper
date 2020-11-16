import { App } from 'vue';

export const install: (app: App) => any;

export class HoperComponent {
    static name: string;

    static install: (app: App) => any;

    $props: Record<string, any>;
}


export class HoperPlugin extends HoperComponent {}