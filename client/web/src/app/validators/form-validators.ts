import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

const USERNAME_MIN = 3;
const USERNAME_MAX = 100;
const PASSWORD_MIN = 3;
const PASSWORD_MAX = 100;
const USERNAME_PATTERN = /^[a-zA-Z0-9_-]+$/;

export class FormValidators {
  static username(): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const value = control.value;
      if (!value) return null;

      if (value.length < USERNAME_MIN || value.length > USERNAME_MAX) {
        return { username: `Username must be between ${USERNAME_MIN} and ${USERNAME_MAX} characters` };
      }
      if (!USERNAME_PATTERN.test(value)) {
        return {
          username:
            'Username must only contain lowercase and uppercase letters, numbers, hyphens, and underscores',
        };
      }
      return null;
    };
  }

  static password(): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const value = control.value;
      if (!value) return null;

      if (value.length < PASSWORD_MIN || value.length > PASSWORD_MAX) {
        return { password: `Password must be between ${PASSWORD_MIN} and ${PASSWORD_MAX} characters` };
      }
      return null;
    };
  }

  static matchField(fieldName: string): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const matchingControl = control.parent?.get(fieldName);
      if (!matchingControl) return null;
      if (control.value !== matchingControl.value) {
        return { mismatch: 'Passwords do not match' };
      }
      return null;
    };
  }
}
