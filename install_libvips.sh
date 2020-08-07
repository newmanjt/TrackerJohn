sudo apt-get install automake build-essential git gobject-introspection \
  libglib2.0-dev libjpeg-turbo8-dev libpng12-dev gtk-doc-tools
git clone https://github.com/jcupitt/libvips.git
cd libvips
./bootstrap.sh
./configure --enable-debug=no --without-python --without-fftw --without-libexif \
  --without-libgf --without-little-cms --without-orc --without-pango --prefix=/usr
make
sudo make install
sudo ldconfig
