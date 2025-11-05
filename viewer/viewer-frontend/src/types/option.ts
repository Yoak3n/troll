import type { DropdownOption } from "naive-ui";

export interface DropdownOptionProps{
    x: number;
    y: number;
    options?: DropdownOption[];
    showDropdown: boolean;
    OnClickoutside: () => void;
    HandleSelect: (key: string) => void
}