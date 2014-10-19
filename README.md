# php-classname-fixer

`php-classname-fixer` is a tool to help rewrite PHP class names.

It can be used to convert code to use [PHP namespaces](http://php.net/manual/en/language.namespaces.php) or rename classes.

Note: `php-classname-fixer` uses regexes to rewrite code. For that reason, it is assumed that files use the [PSR-2](http://www.php-fig.org/psr/psr-2/) code style. It is recommended that your code has been formatted that way before using this tool.

## How it works

`php-classname-fixer` searches a directory for php files. It builds a map of classnames, then determines the implied PSR-0 classname from the location of the files. It then renames classes to use the implied PSR-0 class name.

So for example, if I have a class called `MyClass.php`
```php
class OldVendor_OldNamespace_MyClass {
    public function getDate() {
        return DateTime("2014-10-20");
    }
}
```

and I put that class in `src/NewVendor/NewNamespace/MyClass` and then run `php-classname-fixer src`, the file contents now look like this:

```php
namespace NewVendor\NewNamespace;

class MyClass {
    public function getDate() {
        return \DateTime("2014-10-20");
    }
}
```

In addition all references to the old classname found in the source directory are rewritten.
