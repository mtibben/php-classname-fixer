# php-classname-fixer

`php-classname-fixer` is a tool to help rewrite PHP class names.

It can be used to convert code to use [PHP namespaces](http://php.net/manual/en/language.namespaces.php) or rename classes.

## How it works

`php-classname-fixer` searches a directory for php files. It builds a map of classnames, then determines the implied PSR-0 classname from the location of the files. It then renames classes to use the implied PSR-0 class name.

So for example, if I have a class called `MyClass.php`
```php
class OldVendor_OldNamespace_MyClass {
    
}
```

and I put that class in `src/NewVendor/NewNamespace/MyClass` and then run `php-classname-fixer src`, the file contents now look like this:

```php
namespace NewVendor\NewNamespace;

class MyClass {
    
}
```
